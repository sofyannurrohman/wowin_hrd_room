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
RUN echo "server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
        try_files \$uri \$uri/ /index.html; \
    } \
}" > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
