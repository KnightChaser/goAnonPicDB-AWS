FROM golang:1.21-alpine3.18 as gobuilder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o access_client main.go

FROM alpine:3.19.1
COPY .env ./
RUN export $(cat .env | xargs)
COPY --from=gobuilder /build .
EXPOSE ${CLIENT_WEB_ACCESS_PORT}

CMD ["./access_client"]