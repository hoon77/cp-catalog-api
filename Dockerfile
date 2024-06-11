FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod go.sum main.go ./
COPY common ./common
COPY config  ./config
COPY docs ./docs
COPY handler ./handler
COPY middleware ./middleware
COPY router ./router

RUN go mod download
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

FROM alpine
COPY --from=builder /dist/main .
COPY localize ./localize
COPY repo-cacerts ./repo-cacerts
COPY app.env .
RUN mkdir -p /opt/helm/config
RUN mkdir -p /opt/helm/cache
COPY repositories.yaml /opt/helm/config
ENTRYPOINT ["/main"]