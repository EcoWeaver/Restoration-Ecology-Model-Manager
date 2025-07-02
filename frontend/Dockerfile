FROM node:alpine AS build

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm ci

COPY src src
COPY static static
COPY vite.config.ts .
COPY svelte.config.js .

RUN npm run build

WORKDIR /app/build

COPY package.json .
COPY package-lock.json .

RUN npm ci --omit dev

FROM node:alpine

COPY --from=build /app/build /build

EXPOSE 3000

CMD ["node", "/build"]
