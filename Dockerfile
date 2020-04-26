FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY bin/custm-chat /data/apps/custm-chat/bin/
COPY webim/conf/config.toml /data/apps/custm-chat/conf/
COPY webim/db/ipiptest.ipdb /data/apps/custm-chat/conf/

EXPOSE 8090

ENTRYPOINT ["/data/apps/custm-chat/bin/custm-chat"]
CMD ["-config", "/data/apps/custm-chat/conf/config.toml"]