FROM golang:alpine

RUN apk --no-cache add curl
ENV PGSSLROOTCERT=/.postgresql/root.crt

WORKDIR /app
COPY go.* ./
RUN go mod download

RUN mkdir .postgresql
COPY db/root.crt .postgresql/root.crt

RUN ls .postgresql

EXPOSE 8080

COPY . ./
RUN go build -v -o twitter-bot

CMD ["/app/twitter-bot"]

