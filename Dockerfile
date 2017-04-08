# Docker image for the google container registry plugin
#
#     docker build --rm=true -t plugins/drone-gcr .
FROM rancher/docker:1.10.0
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# ADD drone-gcr /bin/
ADD drone-gcr /go/bin/
VOLUME /var/lib/docker
# ENTRYPOINT ["/bin/drone-gcr"]
ENTRYPOINT ["/go/bin/drone-gcr"]
