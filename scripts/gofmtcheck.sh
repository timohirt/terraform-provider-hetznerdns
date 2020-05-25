#!/usr/bin/env bash
echo "==> Checking that code complies with gofmt requirements..."
files=$(gofmt -l `find . -name '*.go'`)
if [[ -n ${files} ]]; then
    echo 'The following files do not comply with gofmt requirements:'
    echo "${files}"
    echo "Run \`make fmt\` to reformat your code."
    exit 1
fi

exit 