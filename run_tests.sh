#!/bin/bash -e 

[ -d ~/go/bin ] && export PATH="$HOME/go/bin:$PATH"

cd "`dirname \"$0\"`"

DIRS="./cmd/* ./internal/* ./pkg/*"

set -x
go test -cover $DIRS
go vet $DIRS
golint $DIRS

set +x
for cmd in ./cmd/*; do
	set -x
	go build -o /dev/null "$cmd"
	set +x
done

