FROM golang:1.20.2-alpine
WORKDIR /opt
ADD .  /opt

ENV  GOPROXY=https://goproxy.cn,direct 

ENV  DEBUG false

RUN go build -o anicat 

EXPOSE 8080

CMD ["/bin/sh", "-c", "/opt/anicat -d $DEBUG"]
