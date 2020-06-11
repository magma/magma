FROM node:12

EXPOSE 8080

# Create app directory
WORKDIR /app

RUN curl https://get.wasmer.io -sSfL | sh

# Copy package dependencies
COPY package.json yarn.lock babel.config.js ./
COPY fbcnms-projects/workflows-wasm-workers/package.json fbcnms-projects/workflows-wasm-workers/

# Install shared dependencies
COPY fbcnms-packages fbcnms-packages
RUN yarn install --frozen-lockfile && yarn cache clean

# Copy app source
WORKDIR /app/fbcnms-projects/workflows-wasm-workers
COPY fbcnms-projects/workflows-wasm-workers .

# Start app
CMD yarn start
