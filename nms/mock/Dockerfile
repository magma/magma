FROM node:16.14-alpine as builder

WORKDIR /usr/src/

# Copy project dependencies
COPY package.json yarn.lock babel.config.js ./

# Install node dependencies
RUN yarn install --mutex network --frozen-lockfile && yarn cache clean

# Build our static files
COPY . .

FROM node:16.14-alpine

# Install required binaries
RUN apk add ca-certificates curl bash
COPY wait-for-it.sh /usr/local/bin

COPY --from=builder /usr/src /usr/src

WORKDIR /usr/src
CMD ["yarn run mockserver"]
