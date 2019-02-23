FROM node:10-alpine as builder

RUN apk add python g++ make libx11 glew-dev libxi-dev ca-certificates

WORKDIR /usr/src/

# Copy project dependencies
COPY package.json yarn.lock babel.config.js ./
COPY fbcnms-projects/magmalte/package.json fbcnms-projects/magmalte/
COPY fbcnms-packages fbcnms-packages

# Install node dependencies
RUN yarn install --frozen-lockfile && yarn cache clean

# Build our static files
COPY fbcnms-projects/magmalte fbcnms-projects/magmalte
WORKDIR /usr/src/fbcnms-projects/magmalte
RUN yarn run build

FROM node:10-alpine

COPY --from=builder /usr/src /usr/src

# Install required binaries
RUN apk add ca-certificates curl

# Copy startup scripts
COPY fbcnms-projects/magmalte/docker/bin /usr/local/bin
ADD https://github.com/segmentio/chamber/releases/download/v2.2.0/chamber-v2.2.0-linux-amd64 /usr/local/bin/chamber
RUN chmod +x /usr/local/bin/chamber

WORKDIR /usr/src/fbcnms-projects/magmalte
CMD ["start"]
