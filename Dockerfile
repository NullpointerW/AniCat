FROM golang:1.20.2-alpine
WORKDIR /opt
ADD .  /opt

ENV  GOPROXY=https://goproxy.cn,direct 

RUN go build -o anicat 

EXPOSE 8080

CMD ["/opt/anicat"]