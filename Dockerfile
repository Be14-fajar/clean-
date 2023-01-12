FROM golang:1.19.4-alpine as build

# membuat direktori app
RUN mkdir /app

# set working directory /app
WORKDIR /app

COPY ./app

RUN go mod tidy

RUN go build -o belajar

CMD ["/app/belajar"]
