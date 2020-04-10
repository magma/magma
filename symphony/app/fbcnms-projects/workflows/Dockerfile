FROM node:12

EXPOSE 8080

# # Copy the app
# WORKDIR /app
# COPY package.json ./src/ ./

# # Install shared dependencies
# COPY fbcnms-packages fbcnms-packages
# RUN yarn install --frozen-lockfile && yarn cache clean

# RUN yarn install
# CMD yarn start

# Create app directory
WORKDIR /app

# Copy package dependencies
COPY package.json yarn.lock babel.config.js ./
COPY fbcnms-projects/workflows/package.json fbcnms-projects/workflows/

# Install shared dependencies
COPY fbcnms-packages fbcnms-packages
RUN yarn install --frozen-lockfile && yarn cache clean

# Copy app source
WORKDIR /app/fbcnms-projects/workflows
COPY fbcnms-projects/workflows .

# Start app
CMD yarn start
