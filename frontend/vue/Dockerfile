# Stage 1 - Build
FROM node:20.9 as build-deps
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm cache clean --force
RUN npm install
COPY . ./
RUN npm run build

# Stage 2 - Serve
FROM nginx:1.17-alpine
COPY --from=build-deps /usr/src/app/dist /usr/share/nginx/html
EXPOSE 80
COPY deployment/nginx/nginx.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]



