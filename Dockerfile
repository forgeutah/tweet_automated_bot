FROM golang:1.18-alpine

RUN apk update
RUN apk add git

WORKDIR /app
COPY go.* ./
RUN go mod download

EXPOSE 8080

COPY . ./
ADD /dev/.config.json config.json
RUN chmod +xr run.sh
RUN go build -v -o twitter-bot

CMD ["/app/run.sh"]

