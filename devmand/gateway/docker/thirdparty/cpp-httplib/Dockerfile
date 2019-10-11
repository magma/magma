FROM debian:stretch AS %%PKG%%

ARG PKG_DIR=/cache/%%PKG%%
ARG PKG_REPO_DIR=/cache/%%PKG%%/repo
ARG PKG_INSTALL_DIR=/cache/%%PKG%%/install

RUN %%INSTALL%% git %%DEPS%%

WORKDIR $PKG_DIR
RUN git clone %%URL%% $PKG_REPO_DIR

WORKDIR $PKG_REPO_DIR
RUN git checkout %%VERSION%%

RUN mkdir -p $PKG_INSTALL_DIR/include
RUN cp httplib.h $PKG_INSTALL_DIR/include
