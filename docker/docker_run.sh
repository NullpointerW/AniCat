#!/bin/bash

# command-args
qbt_webui_port=$1
qbt_config=$2
usr_bgmi_path=$3
anicat_port=$4
anicat_cfg=$5

# docker run
docker network create anicat-net

docker run -d --name=qb -p ${qbt_webui_port}:8989 \
 -e PUID=1000 -e PGID=1000 -e TZ=Asia/Shanghai -e WEBUIPORT=8989 \
 -v ${qbt_config}:/config -v ${usr_bgmi_path}:/bangumi\
 --restart unless-stopped \
 --network anicat-net --network-alias qb \
 superng6/qbittorrentee:latest

docker run -d --name=anicat --restart unless-stopped \
 -v ${usr_bgmi_path}:/bangumi  -v ${anicat_cfg}:/opt/env.yaml \
 -p ${anicat_port}:8080  --user 1000:1000 \
 --network anicat-net --network-alias anicat \
 wmooon/anicat:latest 