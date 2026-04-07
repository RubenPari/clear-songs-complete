# Build stage
FROM node:20-alpine AS build
WORKDIR /app

# Add ARG for build-time environment variables
ARG API_URL
ARG SPOTIFY_CLIENT_ID
ARG SPOTIFY_REDIRECT_URI

# Set ENV from ARG so generate-env.js can pick them up
ENV API_URL=$API_URL
ENV SPOTIFY_CLIENT_ID=$SPOTIFY_CLIENT_ID
ENV SPOTIFY_REDIRECT_URI=$SPOTIFY_REDIRECT_URI

COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps
COPY . .

# Generate typed environment and build
RUN npm run env && npm run build

# Production stage
FROM nginx:alpine
# In Angular 19/20 with 'application' builder, output is in dist/clear-songs-front/browser
COPY --from=build /app/dist/clear-songs-front/browser /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
