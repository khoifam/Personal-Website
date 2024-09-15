FROM golang:1.23.0-bullseye AS base
WORKDIR /usr/src/app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify
COPY backend/ .
RUN go build -o api

FROM node:12.22.9
WORKDIR /usr/src/app/weather_app
COPY weather_app/package.json weather_app/package-lock.json ./
RUN npm install
COPY weather_app/ .
RUN npm run build

WORKDIR /usr/src/app
COPY --from=base /usr/src/app/backend backend/

WORKDIR /usr/src/app/backend

EXPOSE 8080

CMD ["./api"]