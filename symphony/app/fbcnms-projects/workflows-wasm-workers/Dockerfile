FROM node:12

EXPOSE 8080

# Create app directory
WORKDIR /app

RUN curl https://get.wasmer.io -sSfL | sh

RUN git clone https://github.com/wapm-packages/python.git

# Copy package dependencies
COPY package.json yarn.lock babel.config.js ./
COPY fbcnms-projects/workflows-wasm-workers/package.json fbcnms-projects/workflows-wasm-workers/

# Install shared dependencies
COPY fbcnms-packages fbcnms-packages
RUN yarn install --frozen-lockfile && yarn cache clean

# Copy app source
WORKDIR /app/fbcnms-projects/workflows-wasm-workers
COPY fbcnms-projects/workflows-wasm-workers .

RUN cp -r /app/python/wapm/lib wasm/python/lib
RUN cp wasm/python/bin/python.wasm wasm/python/lib/
# get rid of warning: Could not find platform dependent libraries <exec_prefix>
RUN mkdir wasm/python/lib/python3.6/lib-dynload/

# Start app
CMD yarn start
