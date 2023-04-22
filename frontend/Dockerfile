FROM node:18-slim AS builder

WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile
COPY . ./
RUN yarn build

CMD [ "yarn", "start" ]
