ARG GO_VERSION=1.17.5

FROM golang:${GO_VERSION} as builder

RUN apt-get update

RUN apt-get install -y make

WORKDIR /go/src/game-services

ADD ./ ./

RUN make binary

FROM alpine

RUN apk --no-cache add  \
    ca-certificates \
    libc6-compat \
    libstdc++ \
    file \
    tzdata

RUN mkdir /opt/services

COPY --from=builder /go/src/game-services/bin/* /opt/services/

RUN chmod +x /opt/services

WORKDIR /opt/services