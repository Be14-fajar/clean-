FROM golang:1.19-alpine

##buat folder APP
RUN mkdir /app

##set direktori utama
WORKDIR /app

##copy seluruh file ke app
ADD . /app

##buat executeable
RUN go build -o belajar .

##jalankan executeable
CMD ["/app/belajar"]
