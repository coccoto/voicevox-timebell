FROM golang:1.24.4-bookworm

RUN apt update && apt full-upgrade -y && apt autoremove -y && apt clean
RUN apt install -y alsa-utils

WORKDIR /app

COPY ./server /app/server
COPY ./client /app/client

RUN mkdir -p /app/storage

RUN go build -C /app/server -o /app/dist/app .

CMD ["/app/dist/app"]
