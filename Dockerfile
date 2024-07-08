# Base image with Node.js and Alpine Linux
FROM node:lts-alpine AS base

# Set working directory
WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Build the Next.js application
COPY . .
RUN npm run build

# Install ffmpeg for your application
RUN apk add --no-cache ffmpeg

# Production image
FROM node:lts-alpine AS runner

# Set working directory
WORKDIR /app

# Set environment to production
ENV NODE_ENV production

# Create a system user and group for running the app
RUN addgroup -S nextjs && adduser -S nextjs -G nextjs

# Copy the built application and necessary files
COPY --from=base /app/public ./public
COPY --from=base /app/.next/standalone ./
COPY --from=base /app/.next/static ./.next/static

# Set correct permissions
RUN chown -R nextjs:nextjs /app

# Switch to non-root user
USER nextjs

# Expose the application port
EXPOSE 3000

# Set the environment variable for the app port
ENV PORT 3000

# Run the application
CMD ["node", "server.js"]