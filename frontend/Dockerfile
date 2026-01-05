FROM node:alpine AS build

WORKDIR /app

COPY package.json package-lock.json .

RUN npm ci

COPY . .

RUN npm run build
RUN npm prune --production

FROM node:alpine

COPY --from=build /app/build /build
COPY --from=build /app/node_modules /node_modules

EXPOSE 3000

CMD ["node", "/build"]
