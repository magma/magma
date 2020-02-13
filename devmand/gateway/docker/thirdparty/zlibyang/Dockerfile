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

# Adding prefix to all names as not to interfere with libyang bundled with YDK 
RUN /bin/bash -c "find . \( -name '*.c' -o -name '*.h' -o -name '*.in' \) -exec sed -i -e 's/\bly/llly/g' {} -e 's/\bLY/LLLY/g' -e 's/\bLIBYANG_EXTENSIONS_PLUGINS_DIR/LLLIBYANG_EXTENSIONS_PLUGINS_DIR/g' {} \; "

WORKDIR $PKG_REPO_DIR/build
RUN cmake -DENABLE_LYD_PRIV=ON \
          -DENABLE_CACHE=OFF \
          -DCMAKE_INSTALL_PREFIX=$PKG_INSTALL_DIR \
          -DCMAKE_BUILD_TYPE:String="Release" .. 
RUN make
RUN make install
