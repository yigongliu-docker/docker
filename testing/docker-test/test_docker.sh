#!/bin/bash
set -x
COMMIT=${1-HEAD}
REPO=${2-http://github.com/dotcloud/docker}
BRANCH=${3-master}

# Properly mount cgroups in the container
mount -t tmpfs none /tmp
mount -t tmpfs none /sys/fs/cgroup
cd /sys/fs/cgroup
for C in $(awk "{print \$1}" < /proc/cgroups | grep -v subsys | grep -v memory) ; do mkdir $C ; mount -t cgroup none -o $C $C ; done
pushd /proc/self/fd; for FD in *; do case "$FD" in [012]) ;; *) eval exec "$FD>&-" ;; esac done; popd

# Get docker and its dependencies
export GOPATH=/go
rm -rf $GOPATH
mkdir -p $GOPATH
go get -d github.com/dotcloud/docker
cd /go/src/github.com/kr/pty;          git reset --hard  27435c699
cd /go/src/github.com/gorilla/context; git reset --hard  708054d61e5
cd /go/src/github.com/gorilla/mux;     git reset --hard  9b36453141c
cd /go/src/github.com/dotcloud/tar;    git reset --hard  e5ea6bb21a3294
cd /go/src/code.google.com/p/go.net;   hg checkout -r 84a4013f96e0
cd /go/src/code.google.com/p; hg clone http://code.google.com/p/getopt; cd getopt; hg checkout -r b23bed28ee5c

# Rebase commit in top of master
cd /go/src/github.com/dotcloud/docker
git fetch "$REPO" "$BRANCH:CKT-$BRANCH"
git reset --hard "$COMMIT"
git rebase master

# Test commit
go test -v
