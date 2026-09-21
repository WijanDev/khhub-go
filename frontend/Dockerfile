FROM node:24-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# API origin is injected at container start (see docker-entrypoint.sh).
ARG VITE_API_URL=
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM ghcr.io/static-web-server/static-web-server:2-alpine
USER root
COPY --from=build /app/dist /public
COPY static-web-server.toml /etc/config.toml
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENV SERVER_CONFIG_FILE=/etc/config.toml
ENV SERVER_FALLBACK_PAGE=/public/index.html
# Overridden per environment. Default keeps existing production hosts working.
ENV KHHUB_API_URL=https://api.khhub.app
EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
