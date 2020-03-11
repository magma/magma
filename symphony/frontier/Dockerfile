# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Multi stage application builder
FROM golang:1.14-alpine AS builder

# Install build tools
RUN apk add --no-cache git gcc musl-dev

# Download dependencies
WORKDIR /src/symphony
COPY go.* ./
RUN go mod download

# Copy shared components
COPY pkg /src/symphony/pkg

# Build service
WORKDIR /src/symphony/frontier
COPY frontier .
RUN go build -o /bin/frontier

# Create deployable image
FROM alpine

# Install runtime requirements
RUN apk add --no-cache musl curl

# Copy application artifacts
COPY --from=builder /bin/frontier /bin/

# Set exposed port and entrypoint
EXPOSE 80
ENTRYPOINT ["/bin/frontier"]
