FROM golang:1.13 as build
WORKDIR /go/src/github.com/jukeizu/sentiment
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "sentiment:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/jukeizu/sentiment/passwd /etc/passwd
COPY --from=build --chown=100:101 /go/src/github.com/jukeizu/sentiment/bin/sentiment .
USER nobody
ENTRYPOINT ["./sentiment"]
