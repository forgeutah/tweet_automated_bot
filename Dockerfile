FROM golang:alpine

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o autobot

CMD ["/app/autobot"]

