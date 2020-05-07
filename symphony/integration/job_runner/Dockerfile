FROM networld/grpcurl

RUN apk update \
 && apk add jq \
 && rm -rf /var/cache/apk/*

COPY integration/job_runner/*.sh ./

ENTRYPOINT ["/bin/sh"]
