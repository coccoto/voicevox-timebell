FROM golang:1.24.4-bookworm

WORKDIR /app

COPY ./server /app/server
COPY ./client /app/client

RUN mkdir /app/dist
RUN go build -C /app/server -o /app/dist/app .

CMD ["/app/dist/app"]