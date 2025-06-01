FROM golang:1.24-alpine as build

WORKDIR /go/src/github.com/InsideGallery/brf.im
ADD . .
RUN go build -o /go/bin/server ./cmd/server
CMD ["/go/bin/server"]

FROM alpine:3.7

COPY  --from=build /go/bin/server /go/bin/server
CMD ["/go/bin/server"]
