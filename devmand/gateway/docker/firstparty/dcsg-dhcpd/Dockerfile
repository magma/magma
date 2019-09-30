FROM debian:stretch AS %%PKG%%

RUN %%INSTALL%% %%DEPS%%
RUN mkdir -p /var/lib/dhcp/
RUN touch /var/lib/dhcp/dhcpd.leases

COPY dcsg-dhcpd/dhcpd.conf /dhcpd.conf
COPY dcsg-dhcpd/restart_dhcpd.sh /bin/restart_dhcpd.sh
COPY dcsg-dhcpd/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
