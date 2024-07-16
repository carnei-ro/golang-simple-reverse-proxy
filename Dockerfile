### Builder
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/carnei-ro/golang-simple-reverse-proxy

COPY . ./

RUN go get -d -v ./... && \
  CGO_ENABLED=0 \
  GO111MODULE=on \
  GOOS=${TARGETOS} \
  GOARCH=${TARGETARCH} \
  go build -ldflags="-w -s" -o /go/bin/golang-simple-reverse-proxy

### Final stage
FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source=https://github.com/carnei-ro/golang-simple-reverse-proxy

WORKDIR /

COPY --from=builder /go/bin/golang-simple-reverse-proxy /usr/local/bin/golang-simple-reverse-proxy

USER nonroot:nonroot

CMD [ "/usr/local/bin/golang-simple-reverse-proxy" ]
