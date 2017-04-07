# Docker image for the google container registry plugin
#
#     docker build --rm=true -t plugins/drone-gcr .
FROM alpine:3.2
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD drone-gcr /bin/
ENTRYPOINT ["/bin/drone-gcr"]
