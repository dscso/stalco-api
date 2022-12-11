FROM golang:1.17 AS build-env

ADD . /rest-go
WORKDIR /rest-go
RUN make install
RUN make build

# Final stage
FROM debian:buster

EXPOSE 8000

WORKDIR /
COPY --from=build-env /rest-go/build/apiserver /

CMD ["/apiserver"]