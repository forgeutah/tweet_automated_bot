FROM golang:alpine@latest

RUN apk update
RUN apk add git

WORKDIR /app
COPY go.* ./
RUN go mod download

EXPOSE 8080

COPY . ./
RUN go build -v -o twitter-bot

CMD ["/app/twitter-bot"]

