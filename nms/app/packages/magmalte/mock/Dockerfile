FROM node:12.18-alpine as builder

WORKDIR /usr/src/

# Copy project dependencies
COPY package.json yarn.lock babel.config.js ./
COPY packages/magmalte/package.json packages/magmalte/

# Install node dependencies
RUN yarn install --mutex network --frozen-lockfile && yarn cache clean

# Build our static files
COPY packages packages

FROM node:10-alpine

# Install required binaries
RUN apk add ca-certificates curl bash
COPY packages/magmalte/wait-for-it.sh /usr/local/bin

COPY --from=builder /usr/src /usr/src

WORKDIR /usr/src/packages/magmalte
CMD ["yarn run mockserver"]
