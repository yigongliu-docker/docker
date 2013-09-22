# VERSION:        1.2
# DOCKER-VERSION  0.6.2-dev
# AUTHOR:         Daniel Mizyrycki <daniel@dotcloud.com>
# DESCRIPTION:    Testing docker PRs and commits on top of master using Docker in Docker.
# REFERENCES:     This code reuses the excellent implementation of docker in docker
#                 made by Jerome Petazzoni.  https://github.com/jpetazzo/dind
# COMMENTS:
#   This Dockerfile adapts /Dockerfile to enable docker PRs and commits testing
#   Optional arguments:
#       [commit]  (default: 'HEAD')
#       [repo]    (default: 'http://github.com/dotcloud/docker')
#       [branch]  (default: 'master')
# TO_BUILD:       docker build -t test_docker .
# TO_RUN:         docker run -i -t -privileged test_docker test_docker.sh [commit] [repo] [branch]

from docker
maintainer Daniel Mizyrycki <daniel@dotcloud.com>

# Add go to main /usr/bin
run ln -s /usr/local/go/bin/go /usr/bin/go

# Add test_docker.sh
add test_docker.sh /usr/bin/test_docker.sh
run chmod +x /usr/bin/test_docker.sh
