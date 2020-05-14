# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Multi stage application builder
FROM node:alpine AS builder

# Install application dependencies
WORKDIR /app
COPY package*.json ./
RUN npm install

# Build application
COPY . .
RUN npm run build

# Create production image
FROM node:alpine

# Install application dependencies
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# Copy build dist
COPY --from=builder /app/dist ./dist

# Set entry point and cmd
ENTRYPOINT ["npm"]
CMD ["start"]
