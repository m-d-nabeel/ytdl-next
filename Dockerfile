# Use an official Node.js LTS image as the base image
FROM node:lts-alpine

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN apk add --no-cache ffmpeg \
    && npm install

# Copy the entire project directory to the working directory
COPY . .

# Expose the port Next.js uses (usually 3000)
EXPOSE 3000

# Define the command to start the Next.js app
CMD ["npm", "start"]
