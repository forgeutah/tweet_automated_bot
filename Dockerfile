FROM golang:alpine@latest

apk-get add --update --no-cache git

WORKDIR /app
COPY go.* ./
RUN go mod download

EXPOSE 8080

COPY . ./
RUN go build -v -o twitter-bot

CMD ["/app/twitter-bot"]

