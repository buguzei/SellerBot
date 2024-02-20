FROM golang:latest

COPY ./ ./

RUN go build -o ./bin/SellerBot ./cmd

RUN mkdir "config"

CMD ["./bin/SellerBot"]

