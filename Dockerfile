FROM node:lts-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN apk add ffmpeg && npm install

RUN mkdir -p /usr/src/video
RUN mkdir -p /tmp/video
RUN chmod 777 /tmp/video


COPY . .

RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
