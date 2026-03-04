#!/bin/bash
set -e
cd "$(dirname "$0")"
ln -sfn ../../testdata pkg/analyzer/testdata
echo "Symlink created: pkg/analyzer/testdata -> ../../testdata"
