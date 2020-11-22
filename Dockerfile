FROM golang:1.14 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build .


FROM ubuntu:18.04 AS release

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y ffmpeg

USER root

RUN echo "Europe/Moscow" > /etc/timezone
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow
RUN apt-get install -y tzdata

EXPOSE 5000

COPY --from=build /opt/app/MasterHubBackend /usr/bin/

CMD MasterHubBackend