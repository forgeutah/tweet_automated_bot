FROM alpine as download

RUN apk add curl unzip

ENV ARCH amd64
ENV OP_CLI_VERSION v2.0.0

RUN curl -sSfo op.zip https://cache.agilebits.com/dist/1P/op2/pkg/${OP_CLI_VERSION}/op_linux_${ARCH}_${OP_CLI_VERSION}.zip \
  && unzip -od /usr/local/bin/ op.zip \
  && rm op.zip

FROM alpine

RUN addgroup -S opgroup && adduser -S opuser -G opgroup

RUN apk add libc6-compat
COPY --from=download /usr/local/bin/op /usr/local/bin/op

USER opuser

ENTRYPOINT ["op"]

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

