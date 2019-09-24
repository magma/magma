FROM node:10 AS client

WORKDIR /usr/src/app
COPY . .
RUN npm install

CMD [ "node", "index.js" ]
