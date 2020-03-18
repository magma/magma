FROM node:12.5.0-alpine AS builder

# Install dependencies required by some node modules
RUN apk add --no-cache curl python g++ make libx11 glew-dev libxi-dev

# Copy app source
COPY package.json babel.config.js /app/
COPY fbcnms-packages /app/fbcnms-packages
COPY fbcnms-projects/inventory /app/fbcnms-projects/inventory
COPY fbcnms-projects/magmalte /app/fbcnms-projects/magmalte

# Build static files
WORKDIR /app/fbcnms-projects/inventory
RUN yarn install --frozen-lockfile && yarn cache clean
RUN yarn build

FROM nginx:alpine

RUN apk add --no-cache curl

# Copy static files from build container
COPY --from=builder /app/fbcnms-projects/inventory/static/dist /var/www/inventory/static/dist
COPY fbcnms-projects/inventory/docker/index.html /var/www/

# Copy configs
COPY fbcnms-projects/inventory/docker/nginx.conf /etc/nginx/conf.d/default.conf
