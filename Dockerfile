FROM golang:latest

WORKDIR /app

ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o tigerbeetle_api .

ENTRYPOINT ./tigerbeetle_api