FROM golang:alpine as builder
WORKDIR /loadtest
RUN apk add --no-cache git
COPY . /loadtest
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine
COPY --from=builder /loadtest/weavesocks-loadtest /
ENTRYPOINT ["/weavesocks-loadtest"]