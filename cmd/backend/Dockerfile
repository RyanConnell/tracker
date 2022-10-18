FROM --platform=$BUILDPLATFORM golang:alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /src

ENV CGO_ENABLED=0

RUN apk --update add ca-certificates git

COPY . .

RUN go mod download

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags '-s -w' -o /bin/tracker-backend ./cmd/backend

################################################################################

FROM scratch AS final

USER nobody:nobody

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/tracker-backend /usr/local/bin/tracker-backend

ENTRYPOINT ["/usr/local/bin/tracker-backend"]
