FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM ghcr.io/static-web-server/static-web-server:2
COPY --from=build /app/dist /public
COPY static-web-server.toml /etc/config.toml
ENV SERVER_CONFIG_FILE=/etc/config.toml
ENV SERVER_FALLBACK_PAGE=/public/index.html
EXPOSE 80
