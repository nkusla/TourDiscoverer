FROM node:22.17.0-alpine

WORKDIR /usr/api-gateway

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE ${PORT}

CMD ["npm", "run", "start"]