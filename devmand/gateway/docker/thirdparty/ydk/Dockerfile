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
RUN pip install -r $PKG_REPO_DIR/requirements.txt

# Generate core cpp
RUN $PKG_REPO_DIR/generate.py -v --core --cpp --output-directory=$PKG_BUILD_DIR

# Print generated folders
RUN ls -a $PKG_BUILD_DIR/cpp
RUN ls -a $PKG_BUILD_DIR/cpp/ydk

# Compile core cpp
WORKDIR $PKG_BUILD_DIR/cpp/ydk/build
RUN cmake -DCMAKE_BUILD_TYPE=Release  -DCMAKE_INSTALL_PREFIX=$PKG_INSTALL_DIR ..
RUN make

# Install core cpp
RUN make install

# Define openconfig 2.4.3 bundle file for YDK 
RUN echo '\
{ \
  "name": "openconfig", \
  "version": "2.4.3", \
  "core_version": "0.8.3", \
  "author": "Openconfig", \
  "copyright": "Openconfig", \
  "description": "YDK bundle for OpenConfig models", \
  "models": { \
    "git": [ \
      { \
        "url": "https://github.com/openconfig/public.git", \
        "commits": [ \
          { \
            "commitid": "23c269b2cc6d9d59e4fd5611cd55a9a0c4af3ece", \
            "dir": [ \
              "release/models" \ 
            ] \
          } \
        ] \
      }, \
      { \
        "url": "https://github.com/YangModels/yang.git", \
        "commits": [ \
          { \
            "commitid": "c04cd70b25973e93f51c9bafe764d5219c2a6b40", \
            "file": [ \
              "standard/ietf/RFC/ietf-interfaces@2014-05-08.yang", \
              "standard/ietf/RFC/ietf-inet-types@2013-07-15.yang", \
              "standard/ietf/RFC/ietf-yang-types@2013-07-15.yang", \
              "standard/ietf/RFC/iana-if-type@2017-01-19.yang" \
            ] \
          } \
        ] \
      } \
    ] \
  } \
} \
' > $PKG_REPO_DIR/profiles/bundles/openconfig_pure_2_4_3.json

# Generate openconfig classes
WORKDIR $PKG_REPO_DIR
RUN $PKG_REPO_DIR/generate.py -v --bundle $PKG_REPO_DIR/profiles/bundles/openconfig_pure_2_4_3.json --cpp --output-directory=$PKG_BUILD_DIR || true

# Compile openconfig
WORKDIR $PKG_BUILD_DIR/cpp/openconfig-bundle/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=$PKG_INSTALL_DIR ..
RUN make

# Install openconfig
RUN make install

