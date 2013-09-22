#!/bin/sh -e
# This is a one time script to prepare docker-ci

# Build docker container
cd /go/src/github.com/dotcloud/docker; docker build -t docker .

# Build docker-test container
cd testing/docker-test; docker build -t test_docker .

# Build docker nightly release container
cd ../nightlyrelease; docker build -t dockerbuilder .

# Self removing
echo -e '#!/bin/sh -e\nexit 0\n' > /etc/rc.local
exit 0
