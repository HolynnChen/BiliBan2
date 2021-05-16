FROM alpine:latest
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" >/etc/timezone && apk del tzdata
RUN mkdir -p /data/biliban
COPY ./env_pro.toml /data/biliban/env.toml
COPY ./build/main /data/biliban/
RUN chmod +X /data/biliban/main
WORKDIR /data/biliban
ENTRYPOINT ["/data/biliban/main"]
