FROM golang:alpine

WORKDIR /app
COPY go.* ./
RUN go mod download

EXPOSE 8080

COPY . ./
RUN go build -v -o twitter-bot

CMD ["/app/twitter-bot"]

