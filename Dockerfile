FROM node:lts-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN apk add ffmpeg && npm install

RUN mkdir -p /tmp/downloaded
RUN chmod 777 /tmp/downloaded
COPY . .

RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
