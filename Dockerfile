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
COPY repo-cacerts ./repo-cacerts
COPY router ./router

RUN go mod download
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

FROM scratch
COPY --from=builder /dist/main .
COPY localize ./localize
COPY app.env .
ENTRYPOINT ["/main"]