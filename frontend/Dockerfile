# Stage 1: Build the Vue application
FROM node:20-alpine as builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Stage 2: Serve the application with Nginx
FROM nginx:alpine

# The Nginx config in the root of the project acts as a reverse proxy
# But here we just need a simple static server that handles Vue router history mode
COPY --from=builder /app/dist /usr/share/nginx/html

# Update default nginx config to fallback to index.html for Vue Router history mode
# index.html must NEVER be cached so the browser always gets the latest chunk hashes
RUN echo "server { \
    listen 80; \
    root /usr/share/nginx/html; \
    index index.html index.htm; \
    \
    # Hashed assets (JS/CSS chunks) — cache forever, they are content-addressed \
    location ~* \\.(?:js|css)\$ { \
        expires 1y; \
        add_header Cache-Control \"public, immutable\"; \
        try_files \$uri =404; \
    } \
    \
    # index.html — never cache so new deployments are picked up immediately \
    location = /index.html { \
        add_header Cache-Control \"no-cache, no-store, must-revalidate\"; \
        add_header Pragma \"no-cache\"; \
        add_header Expires 0; \
        try_files \$uri =404; \
    } \
    \
    # SPA fallback for all other routes \
    location / { \
        try_files \$uri \$uri/ /index.html; \
    } \
}" > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
