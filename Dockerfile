ARG GOLANG_VERSION=1.18-alpine

FROM golang:${GOLANG_VERSION} AS build
WORKDIR /build
COPY . .


RUN go mod vendor

RUN go build -o /bin/mclight -mod=vendor

FROM alpine:latest AS dev

COPY --from=build /bin/mclight /bin/mclight

EXPOSE 50070
ENTRYPOINT ["/bin/mclight"]
CMD []
