FROM golang:latest

RUN apt install git

WORKDIR /app

ADD . .
RUN go build -o tigerbeetle_api .

ENTRYPOINT ./tigerbeetle_api