FROM debian:stretch AS %%PKG%%

RUN %%INSTALL%% %%DEPS%%
RUN %%INSTALL%% procps strace sudo vim gdb curl valgrind

COPY --from=%%DOCKER_REGISTRY%%/thirdparty /cache/install /
COPY --from=%%DOCKER_REGISTRY%%/firstparty /cache/install /

COPY --from=%%DOCKER_REGISTRY%%/firstparty /cache/magma/repo/orc8r/gateway/python /tmp/orc8r
RUN pip3 install /tmp/orc8r
COPY --from=%%DOCKER_REGISTRY%%/firstparty /build/python/gen /usr/local/lib/python3.5/dist-packages/

COPY symphony-agent/files/configs/ /etc/magma
COPY symphony-agent/files/magma_system_prepare /bin/

COPY --from=%%DOCKER_REGISTRY%%/firstparty /cache/magma/repo/fb/config/certs/rootCA.pem /var/opt/magma/certs/rootCA.pem

RUN mkdir -p /var/opt/magma/configs/

COPY install /

COPY --from=%%DOCKER_REGISTRY%%/devmand /usr/share/snmp /usr/share/snmp
COPY --from=%%DOCKER_REGISTRY%%/devmand /var/lib/snmp /var/lib/snmp
COPY --from=%%DOCKER_REGISTRY%%/devmand /etc/snmp /etc/snmp

ENV MIBDIRS /usr/share/mibs

COPY symphony-agent/files/devmand.service /etc/systemd/system/devmand.service
COPY symphony-agent/files/magmad.service /etc/systemd/system/magmad.service

# TODO do we really need these?
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN rm -f /lib/systemd/system/multi-user.target.wants/* \
    /etc/systemd/system/*.wants/* \
    /lib/systemd/system/local-fs.target.wants/* \
    /lib/systemd/system/sockets.target.wants/*udev* \
    /lib/systemd/system/sockets.target.wants/*initctl* \
    /lib/systemd/system/sysinit.target.wants/systemd-tmpfiles-setup* \
    /lib/systemd/system/systemd-update-utmp*

VOLUME ["/sys/fs/cgroup"]

RUN systemctl enable devmand.service
RUN systemctl enable magmad.service

CMD ["/lib/systemd/systemd"]
