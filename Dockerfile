FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache make

COPY . .

ENTRYPOINT ["sh"]
