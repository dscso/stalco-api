FROM golang:1.17 AS build-env

RUN mkdir /rest-go
ADD Makefile /rest-go
ADD go.mod /rest-go
ADD go.sum /rest-go

WORKDIR /rest-go
RUN make install

ADD . /rest-go
RUN make build

# Final stage
FROM debian:buster

EXPOSE 8000

WORKDIR /
COPY --from=build-env /rest-go/build/apiserver /

CMD ["/apiserver"]