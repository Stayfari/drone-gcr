# Docker image for the google container registry plugin
#
#     docker build --rm=true -t plugins/drone-gcr .
FROM docker:17.04.0-dind
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /home/root
ENV HOME="/home/root"
ADD drone-gcr /bin/
ENTRYPOINT ["/bin/drone-gcr"]
