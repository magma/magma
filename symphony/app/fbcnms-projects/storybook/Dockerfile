FROM node:13.8 AS builder

WORKDIR /src/

# Copy package dependencies
COPY package.json yarn.lock babel.config.js ./

# copy storybook
COPY fbcnms-packages /src/fbcnms-packages

# install assets
RUN yarn install --frozen-lockfile && yarn cache clean

WORKDIR /src/fbcnms-packages/fbcnms-ui/

# build project
RUN rm -rf storybook-static && yarn build-storybook -s public

# create reverse proxy image
FROM nginx:alpine

# Install runtime dependencies
RUN apk add --no-cache curl

# copy build artifacts
COPY --from=builder /src/fbcnms-packages/fbcnms-ui/storybook-static /www/data

# copy nginx config
COPY fbcnms-projects/storybook/nginx/default.nginx /etc/nginx/conf.d/default.conf
