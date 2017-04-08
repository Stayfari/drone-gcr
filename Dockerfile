# Docker image for the google container registry plugin
#
#     docker build --rm=true -t plugins/drone-gcr .
FROM docker:17.04.0-dind
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD drone-gcr /bin/
ENTRYPOINT ["/bin/drone-gcr"]
