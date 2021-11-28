FROM golang:1.17.3

COPY . /tmp/tracker
WORKDIR /tmp/tracker

ENV GOOS linux
ENV GOARCH amd64

RUN go build -o /bin/tracker-frontend /tmp/tracker/cmd/frontend/frontend.go
RUN go build -o /bin/tracker-backend /tmp/tracker/cmd/backend/backend.go
RUN go build -o /bin/tracker-scraper /tmp/tracker/cmd/scraper/scraper.go

RUN chmod a+x /bin/tracker-*
