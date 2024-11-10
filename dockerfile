FROM alpine:latest

RUN apk add --no-cache go just

COPY . /app
WORKDIR /app
RUN go build

CMD ["sh"]
