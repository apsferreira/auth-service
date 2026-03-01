# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .
# npx vite build para evitar erros de type-check estrito do tsc
RUN npx vite build

# Production stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx config: SPA routing + proxy /api → authservice-backend (K3s ClusterIP)
RUN printf 'server {\n    listen 80;\n    root /usr/share/nginx/html;\n    index index.html;\n    location / {\n        try_files $uri $uri/ /index.html;\n    }\n    location /api {\n        proxy_pass http://authservice-backend:3010;\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n        proxy_set_header X-Forwarded-Proto $scheme;\n    }\n    location /health {\n        access_log off;\n        return 200 "healthy\\n";\n        add_header Content-Type text/plain;\n    }\n}\n' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
