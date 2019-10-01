FROM debian:stretch AS %%PKG%%

ARG PKG_DIR=/cache/%%PKG%%
ARG PKG_REPO_DIR=/cache/%%PKG%%/repo
ARG PKG_BUILD_DIR=/cache/%%PKG%%/build
ARG PKG_INSTALL_DIR=/cache/%%PKG%%/install

RUN %%INSTALL%% git cmake ninja-build %%DEPS%%

COPY --from=%%DOCKER_REGISTRY%%/thirdparty /cache/install /usr
COPY . $PKG_REPO_DIR

ENV MAGMA_ROOT $PKG_REPO_DIR

ENV PYTHON_BUILD /build/python
ENV PIP_CACHE_HOME ~/.pipcache
ENV PROTO_LIST orc8r_protos devmand_protos

WORKDIR $PKG_REPO_DIR/orc8r/gateway/python
RUN make protos

RUN cp -r /usr/include/prometheus/* /usr/include

WORKDIR $PKG_BUILD_DIR
RUN cmake -GNinja -DCMAKE_BUILD_TYPE=release  \
          -DGRPC_CPP_PLUGIN_PATH=/usr/bin \
          -DCMAKE_PREFIX_PATH=/cache/install \
          $PKG_REPO_DIR/orc8r/gateway/c/common
RUN ninja -C .
RUN mkdir -p ${PKG_INSTALL_DIR}/lib
RUN mkdir -p ${PKG_INSTALL_DIR}/include/orc8r/protos
RUN find . -iname "*.a" -exec cp {} ${PKG_INSTALL_DIR}/lib \;
RUN find $PKG_REPO_DIR -iname "*.h" -exec cp {} ${PKG_INSTALL_DIR}/include \;
RUN cp -r ${PKG_BUILD_DIR}/service303/orc8r/protos ${PKG_INSTALL_DIR}/include/orc8r
