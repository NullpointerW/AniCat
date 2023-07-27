# AniCat

<p align="center">
    <img title="pre" src="static/cmd-pre.gif" alt="" width=75%>
</p>

通过命令行交互的一键追番工具

[![Go Report Card](https://goreportcard.com/badge/github.com/NullpointerW/anicat)](https://goreportcard.com/report/github.com/NullpointerW/anicat) [![Go Report Card](https://img.shields.io/badge/go-v1.20-blue)](https://tip.golang.org/doc/go1.20)
## 功能特点
 * 傻瓜式: 无需配置任何网站的账号或者rss订阅，只需输入想要追踪的番剧名即可一键添加
 * 支持下载后的剧集自动重命名，方便媒体软件进行刮削 e.g:
  ```
  [UHA-WINGS]Bocchi the Rock![05][x264][1080p]>>孤独摇滚！S01E05
  ```
 * 支持字幕组筛选、关键字正则过滤
 * 番剧下载完成后推送服务(目前支持邮件)
 
 ## 部署
 ### linux 
 推荐使用 docker-compose 部署
 * 创建应用目录
``` bash
mkdir -p /usr/opt/anicat && cd /usr/opt/anicat
   ```
 * 下载配置文件并修改
``` bash
mkdir cfg && wget -P ./cfg/ https://raw.githubusercontent.com/NullpointerW/AniCat/master/env.yaml && vim ./cfg/env.yaml
   ```
 * 修改配置文件
```yaml
port: 8080 # 监听端口 docker-compose 部署时无需更改
path: /bangumi # 番剧下载路径 docker-compose 部署时无需更改
drop-dumplicate: on # 若存在相同集数，则删除重复项（建议开启)
qbittorrent:
  url: http://qb:8989 # qbt-api url,在docker-compose环境中,host为qb
  username: admin
  password: adminadmin
  localed: yes # 本机登录，如果qbt开启了本地登录选项，则可不用填写用户名和密码，docker-compose部署则可忽视此项
  timeout: 3500 # qbt-api请求的超时时间，有时任务添加到qbt上，调用api后无法立即响应到数据
  proxy: # qbt代理配置 可选项
    address: remote:7890 # 配置qbt的代理地址
    type: http # 类型可为 http,socks5等 详见qbt wiki
    host-lookup: on # 使用代理查询域名
    peer: on # 使用代理进行对端连接(文件传输)

crawl: # 为爬虫设置代理，可省略
  proxies: # 可设置多个代理进行轮询
    - http://remote:7890
    - http://remote:7891
    - http://remote1:7890

push: # 配置推送服务，如无此需求则可省略
  email:
    host: smtp.qq.com # SMTP 
    port: 25
    username: xxx@qq.com # 发件邮箱
    password: xxx
    skipssl: yes # 跳过ssl,开启此项可能需要变更相应的smtp地址，具体情况依照邮箱运营商
    template: tmp/template.html # 邮件模板地址，若省略则使用内置的模板
```
  
 * 下载docker-compose yaml
``` bash 
wget https://raw.githubusercontent.com/NullpointerW/AniCat/master/Docker/docker-compose.yml
   ```

``` yaml
version: "3.9"
services:

  anicat:
    image: wmooon/anicat:latest
    container_name: anicat
    ports:
      - 8080:8080 # anicat监听端口
    environment:
      - DEBUG=false # 关闭debug模式
    depends_on:
      - qb
    user: "1000:1000"
    volumes:
      - ./cfg/env.yaml:/opt/env.yaml # 配置文件路径
      - ./bangumi:/bangumi # 番剧下载路径，如果改动则需要和qb保持一致
    restart: unless-stopped
 

  qb:
    image: superng6/qbittorrentee:latest
    container_name: qb
    ports:
      - 8989:8989 # webui 端口
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Asia/Shanghai
      - WEBUIPORT=8989 
    volumes:
      - ./qb:/config
      - ./bangumi:/bangumi
    restart: unless-stopped
    
```
 * 运行
``` bash 
docker compose up -d
   ```

### windows
 * [下载qbittorrent](https://www.qbittorrent.org/download.php)（版本≥ v4.1）