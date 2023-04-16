FROM node:18-slim

WORKDIR /app
COPY . .
RUN yarn --frozon-lockfile
RUN yarn build

CMD yarn start
