# Docker image for the google container registry plugin
#
#     docker build --rm=true -t plugins/drone-gcr .
FROM rancher/docker:v1.11.2

ADD drone-gcr /go/bin/
VOLUME /var/lib/docker
ENTRYPOINT ["/usr/bin/dockerlaunch", "/go/bin/drone-gcr"]
