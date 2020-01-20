FROM node:10-alpine as builder

RUN apk add python g++ make libx11 glew-dev libxi-dev ca-certificates

WORKDIR /usr/src/

# Copy project dependencies
COPY package.json yarn.lock babel.config.js ./
COPY fbcnms-projects/magmalte/package.json fbcnms-projects/magmalte/
COPY fbcnms-packages fbcnms-packages

# Install node dependencies
RUN yarn install --mutex network --frozen-lockfile && yarn cache clean

# Build our static files
COPY fbcnms-projects/magmalte fbcnms-projects/magmalte
WORKDIR /usr/src/fbcnms-projects/magmalte
RUN yarn run build

FROM node:10-alpine

# Install required binaries
RUN apk add ca-certificates curl bash
COPY fbcnms-projects/magmalte/wait-for-it.sh /usr/local/bin

COPY --from=builder /usr/src /usr/src

WORKDIR /usr/src/fbcnms-projects/magmalte
CMD ["yarn run start:prod"]
