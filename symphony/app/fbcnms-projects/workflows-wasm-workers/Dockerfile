FROM node:12 as wasmer
WORKDIR /wasmer
RUN wget -q https://github.com/wasmerio/wasmer/releases/download/0.17.0/wasmer-linux-amd64.tar.gz -O wasmer.tar.gz && \
 tar xvfz wasmer.tar.gz

FROM wasienv/wasienv as qjs-wasi
RUN apt-get install -y make gcc git
COPY --from=wasmer /wasmer/bin/wasmer /root/.wasienv/bin/
RUN git clone https://github.com/robot-rumble/qjs-wasi.git /qjs
WORKDIR /qjs/src
RUN git checkout 05d9099
RUN make qjs.wasm

FROM node:12
EXPOSE 8080
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user
RUN mkdir -p /app/fbcnms-projects/workflows-wasm-workers && \
 groupadd -r wasm && \
 useradd --no-log-init --system -g wasm --home /app/fbcnms-projects/workflows-wasm-workers wasm && \
 mkdir /app/fbcnms-projects/workflows-wasm-workers/.wasmer && \
 chown wasm:wasm /app/fbcnms-projects/workflows-wasm-workers/.wasmer && \
 mkdir /app/fbcnms-projects/workflows-wasm-workers/.yarn && \
 chown wasm:wasm /app/fbcnms-projects/workflows-wasm-workers/.yarn


# install wasmer wasm runtime to ~/.wasmer/bin/wasmer
#RUN curl https://get.wasmer.io -sSfL | sh
WORKDIR /app/fbcnms-projects/workflows-wasm-workers/.wasmer
COPY --from=wasmer /wasmer/bin/wasmer bin/

# Create app directory
WORKDIR /app

# Copy package dependencies
COPY package.json yarn.lock babel.config.js ./
COPY fbcnms-projects/workflows-wasm-workers/package.json fbcnms-projects/workflows-wasm-workers/

# Install shared dependencies
COPY fbcnms-packages fbcnms-packages
RUN yarn install --frozen-lockfile && yarn cache clean

# Copy python lib files
RUN git clone https://github.com/wapm-packages/python.git /app/python
WORKDIR /app/python
RUN git checkout 7435ecf395417a225a80e2e41e0c673ad17f5344

# Copy app source
WORKDIR /app/fbcnms-projects/workflows-wasm-workers
COPY fbcnms-projects/workflows-wasm-workers .

# Finish python copying
RUN cp -r /app/python/wapm/lib wasm/python/lib
RUN cp wasm/python/bin/python.wasm wasm/python/lib/
RUN rm -rf /app/python
# get rid of warning: Could not find platform dependent libraries <exec_prefix>
RUN mkdir wasm/python/lib/python3.6/lib-dynload/

COPY --from=qjs-wasi /qjs/src/qjs.wasm wasm/quickjs/quickjs.wasm

USER wasm

# Start app
CMD yarn start
