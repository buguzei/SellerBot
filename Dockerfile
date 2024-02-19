FROM golang:latest

COPY ./ ./

RUN go build -o ./bin/SellerBot ./cmd

CMD ["./bin/SellerBot"]

