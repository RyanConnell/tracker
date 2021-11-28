FROM golang:1.17.3

COPY . /tmp/tracker
WORKDIR /tmp/tracker

RUN mv /tmp/tracker/bin/* /bin/
RUN mv /tmp/tracker/web/templates .
RUN mv /tmp/tracker/web/public .
