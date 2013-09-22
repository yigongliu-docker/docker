#!/bin/bash
set -x
COMMIT=${1-HEAD}
REPO=${2-http://github.com/dotcloud/docker}
BRANCH=${3-master}

DOCKER_PATH=/go/src/github.com/dotcloud/docker

# Fetch latest master
cd /; rm -rf $DOCKER_PATH; mkdir -p $DOCKER_PATH; cd $DOCKER_PATH
git init .
git fetch -q http://github.com/dotcloud/docker master
git reset --hard FETCH_HEAD

echo FIXME. Temporarily add Jerome changeset with proper apparmor handling

git fetch -q http://github.com/jpetazzo/docker escape-apparmor-confinement:escape-apparmor-confinement
git rebase --onto master master escape-apparmor-confinement

# Rebase commit in top of master
git fetch -q "$REPO" "$BRANCH:CKT-$BRANCH"
git reset --hard "$COMMIT"
#git rebase master

#### FIXME. Temporarily rebase on top of Jerome changeset
git rebase escape-apparmor-confinement

# Test commit
go test -v
