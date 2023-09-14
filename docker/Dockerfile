FROM golang:1.21.1-alpine AS builder

ARG VER=x.x.x
WORKDIR /build
ADD .  /build
ENV  GOPROXY=https://goproxy.cn,direct 
RUN echo "version---${VER}"
RUN go build  -ldflags " -X "github.com/NullpointerW/anicat/conf.Ver=${VER}"" -o anicat 

FROM alpine:3.8

ENV  DEBUG false

WORKDIR /opt

COPY --from=builder /build/anicat  /opt/anicat

EXPOSE 8080
RUN chmod +x ./anicat
CMD ["/bin/sh", "-c", "./anicat -d=$DEBUG"]
