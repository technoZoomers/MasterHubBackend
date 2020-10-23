FROM golang:1.14 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build .


FROM ubuntu:18.04 AS release

USER root

EXPOSE 5000

COPY --from=build /opt/app/MasterHubBackend /usr/bin/

CMD MasterHubBackend