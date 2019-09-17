FROM debian:stretch AS %%PKG%%
ARG PKG_DIR=/cache/%%PKG%%
ARG PKG_REPO_DIR=/cache/%%PKG%%/repo
ARG PKG_BUILD_DIR=/cache/%%PKG%%/build
ARG PKG_INSTALL_DIR=/cache/%%PKG%%/install

RUN %%INSTALL%% git cmake make %%DEPS%%

WORKDIR $PKG_DIR
RUN git clone %%URL%% $PKG_REPO_DIR

WORKDIR $PKG_REPO_DIR
RUN git checkout %%VERSION%%
RUN git submodule update --init

WORKDIR $PKG_BUILD_DIR/ydk
RUN cmake -DCMAKE_BUILD_TYPE=release \
          -DCMAKE_INSTALL_PREFIX=$PKG_INSTALL_DIR \
          $PKG_REPO_DIR/core/ydk
RUN make
RUN make install

WORKDIR $PKG_BUILD_DIR/ietf
RUN cmake -DCMAKE_BUILD_TYPE=release \
          -DCMAKE_INSTALL_PREFIX=$PKG_INSTALL_DIR \
          $PKG_REPO_DIR/ietf
RUN make
RUN make install

WORKDIR $PKG_BUILD_DIR/openconfig
RUN cmake -DCMAKE_BUILD_TYPE=release \
          -DCMAKE_INSTALL_PREFIX=$PKG_INSTALL_DIR \
          $PKG_REPO_DIR/openconfig
RUN make
RUN make install
