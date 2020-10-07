FROM node:12.18-alpine as builder

RUN apk add python g++ make libx11 glew-dev libxi-dev ca-certificates

WORKDIR /usr/src/

# Copy project dependencies
COPY package.json yarn.lock babel.config.js ./
COPY packages/magmalte/package.json packages/magmalte/
COPY packages/fbcnms-magma-api/package.json packages/fbcnms-magma-api/

# Install node dependencies
ENV PUPPETEER_SKIP_DOWNLOAD "true"
RUN yarn install --mutex network --frozen-lockfile && yarn cache clean

# Build our static files
COPY packages packages
WORKDIR /usr/src/packages/magmalte
RUN yarn run build

FROM node:10-alpine

# Install required binaries
RUN apk add ca-certificates curl bash
COPY packages/magmalte/wait-for-it.sh /usr/local/bin

COPY --from=builder /usr/src /usr/src

WORKDIR /usr/src/packages/magmalte
CMD ["yarn run start:prod"]
