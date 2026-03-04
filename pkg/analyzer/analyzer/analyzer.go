package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"logcheck-selectel/pkg/rules/rules"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var logMethods = map[string]bool{
	"Debug": true, "Debugf": true, "Debugw": true, "DebugContext": true,
	"Info": true, "Infof": true, "Infow": true, "InfoContext": true,
	"Warn": true, "Warnf": true, "Warnw": true, "WarnContext": true,
	"Error": true, "Errorf": true, "Errorw": true, "ErrorContext": true,
	"Fatal": true, "Fatalf": true, "Fatalw": true,
	"Panic": true, "Panicf": true, "Panicw": true,
	"Log": true, "Logf": true,
	"With": false,
}

var logPackages = map[string]bool{
	"log/slog":        true,
	"go.uber.org/zap": true,
}
var configPath string
var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks that log messages follow conventions: lowercase start, English only, no special chars/emoji, no sensitive data",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", ".logcheck.yml", "path to logcheck configuration file")
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfg, _ := LoadConfig(configPath)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		if !isLogCall(pass, call) {
			return
		}
		msg, msgExpr := extractMessage(pass, call)
		if msg == "" || msgExpr == nil {
			return
		}
		checkAndReport(pass, call, msgExpr, msg, cfg)
	})

	return nil, nil
}

func isLogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	methodName := sel.Sel.Name
	if _, ok := logMethods[methodName]; !ok {
		return false
	}

	switch x := sel.X.(type) {
	case *ast.Ident:
		obj := pass.TypesInfo.ObjectOf(x)
		if obj == nil {
			return false
		}
		if pkgName, ok := obj.(*types.PkgName); ok {
			return logPackages[pkgName.Imported().Path()]
		}
		return isLoggerType(pass, x)
	default:
		typ := pass.TypesInfo.TypeOf(sel.X)
		if typ == nil {
			return false
		}
		return isLoggerTypePath(typ.String())
	}
}
func isLoggerType(pass *analysis.Pass, ident *ast.Ident) bool {
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return false
	}
	typ := obj.Type()
	if typ == nil {
		return false
	}
	return isLoggerTypePath(typ.String())
}

func isLoggerTypePath(typStr string) bool {
	return strings.Contains(typStr, "log/slog") ||
		strings.Contains(typStr, "go.uber.org/zap")
}

func extractMessage(pass *analysis.Pass, call *ast.CallExpr) (string, ast.Expr) {
	if len(call.Args) == 0 {
		return "", nil
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok && strings.HasSuffix(sel.Sel.Name, "Context") && len(call.Args) >= 2 {
		return extractStringValue(pass, call.Args[1])
	}

	return extractStringValue(pass, call.Args[0])
}

func extractStringValue(pass *analysis.Pass, expr ast.Expr) (string, ast.Expr) {
	tv, ok := pass.TypesInfo.Types[expr]
	if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		return constant.StringVal(tv.Value), expr
	}

	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		s := strings.Trim(lit.Value, "\"`")
		return s, expr
	}

	if bin, ok := expr.(*ast.BinaryExpr); ok && bin.Op == token.ADD {
		lhs, lexpr := extractStringValue(pass, bin.X)
		if lhs != "" {
			return lhs, lexpr
		}
	}

	return "", nil
}

func checkAndReport(pass *analysis.Pass, call *ast.CallExpr, msgExpr ast.Expr, msg string, cfg Config) {
	if cfg.Rules.LowercaseStart {
		if diag := rules.CheckLowercaseStart(msg); diag != "" {
			fixed := rules.FixLowercaseStart(msg)
			report(pass, msgExpr, diag, msg, fixed)
		}
	}

	if cfg.Rules.EnglishOnly {
		if diag := rules.CheckEnglishOnly(msg); diag != "" {
			pass.Reportf(msgExpr.Pos(), "%s", diag)
		}
	}

	if cfg.Rules.NoSpecialChars {
		if diag := rules.CheckNoSpecialChars(msg, cfg.AllowedSpecialChars); diag != "" {
			fixed := rules.FixNoSpecialChars(msg)
			report(pass, msgExpr, diag, msg, fixed)
		}
	}

	if cfg.Rules.NoSensitive {
		if diag := rules.CheckNoSensitiveData(msg, cfg.SensitivePatterns); diag != "" {
			pass.Reportf(msgExpr.Pos(), "%s", diag)
		}
	}
}

func report(pass *analysis.Pass, expr ast.Expr, message, oldMsg, newMsg string) {
	if oldMsg == newMsg {
		pass.Reportf(expr.Pos(), "%s", message)
		return
	}

	newText := `"` + newMsg + `"`

	pass.Report(analysis.Diagnostic{
		Pos:     expr.Pos(),
		End:     expr.End(),
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "replace with: " + newText,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(newText),
					},
				},
			},
		},
	})
}
