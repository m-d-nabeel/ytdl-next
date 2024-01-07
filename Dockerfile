FROM node:lts-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN apk add --no-cache ffmpeg \
    && npm install

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
