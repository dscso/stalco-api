FROM golang:1.17 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev
RUN make install
RUN make swag
RUN make build

# Final stage
FROM debian:buster

EXPOSE 8000

WORKDIR /
COPY --from=build-env /dockerdev/build/apiserver /

CMD ["/apiserver"]