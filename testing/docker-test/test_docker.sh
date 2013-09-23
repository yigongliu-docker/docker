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
git pull -q https://github.com/jpetazzo/docker.git escape-apparmor-confinement || exit 1

# Rebase commit in top of master
git fetch -q "$REPO" "$BRANCH"
git merge --no-edit $COMMIT || exit 1

# Test commit
go test -v
