FROM node:8.11.4

WORKDIR /app/website

EXPOSE 3000 35729
COPY readmes /app/docs
COPY docusaurus /app/website
RUN yarn install
RUN yarn build

CMD ["yarn", "start"]
