FROM linuxserver/openssh-server:8.6_p1-r2-ls56
RUN apk add --no-cache gettext

ENV TPL="config.template"

COPY "$TPL" "/$TPL"
COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
