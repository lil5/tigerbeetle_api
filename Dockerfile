FROM golang:latest

WORKDIR /app

ADD . .
RUN go build -o tigerbeetle_api .

ENTRYPOINT ./tigerbeetle_api