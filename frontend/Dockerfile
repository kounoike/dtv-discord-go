FROM node:18-slim AS builder

WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile
COPY . ./
RUN rm public/index && rm public/encoded && rm public/recorded
RUN yarn build

FROM nginx:latest
COPY --from=builder /app/dist /www
