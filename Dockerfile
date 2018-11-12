FROM golang:1.11.1 as build

COPY . /resiproxy
WORKDIR /resiproxy/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:edge
COPY --from=build /go/src/github.com/marqub/resiproxy/resiproxy /resiproxy

EXPOSE 8080
ENTRYPOINT ["/resiproxy"]   