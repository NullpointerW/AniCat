FROM golang:1.20.2-alpine AS builder

WORKDIR /build
ADD .  /build

ENV  GOPROXY=https://goproxy.cn,direct 
RUN go build -o anicat 

FROM alpine:3.8

ENV  DEBUG false

WORKDIR /opt

COPY --from=builder /build/anicat  /opt/anicat

EXPOSE 8080

CMD ["/bin/sh", "-c", "./anicat -d $DEBUG"]
