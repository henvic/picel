FROM golang:onbuild
MAINTAINER Henrique Vicente <henriquevicente@gmail.com>

RUN apt-get update
RUN apt-get install -y \
    imagemagick \
    webp

EXPOSE 8123
