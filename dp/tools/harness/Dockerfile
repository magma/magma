FROM python:3.8.9-alpine3.13 as builder

RUN apk add --no-cache \
        gcc=10.2.1_pre1-r3 \
        musl-dev=1.2.2-r0 \
        libffi-dev=3.3-r2 \
        openssl-dev=1.1.1k-r0 \
	git=2.30.2-r0 \
	rust=1.47.0-r2 \
	cargo=1.47.0-r2 \
	linux-headers=5.7.8-r0 \
	g++=10.2.1_pre1-r3 \
	curl-dev=7.76.1-r0 \
	python3-dev=3.8.10-r0 \
	geos-dev=3.8.1-r2 \
	libressl-dev=3.1.5-r0 \
	libxslt-dev=1.1.34-r0 \
        && pip install --no-cache-dir virtualenv==20.10.0

WORKDIR /tests
ENV CLONE_URL="https://github.com/Wireless-Innovation-Forum/Spectrum-Access-System.git"
RUN git clone \
	--depth 1 \
	--filter=blob:none \
	--sparse $CLONE_URL .; \
	git sparse-checkout set src/harness

RUN virtualenv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
COPY requirements.txt /tests/requirements.txt
RUN pip install --no-cache-dir -r requirements.txt
	

FROM python:3.8.9-alpine3.13
RUN apk add --no-cache \
	geos-dev=3.8.1-r2 \
	libxslt-dev=1.1.34-r0 \
	libcurl=7.76.1-r0
COPY --from=builder /opt/venv /opt/venv
COPY --from=builder /tests/src/harness /opt/server
ENV PATH="/opt/venv/bin:$PATH"
WORKDIR /opt/server
RUN chown -R nobody /opt/server
USER nobody
ENTRYPOINT ["python"]
CMD ["test_main.py"]
