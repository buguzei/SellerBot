FROM golang:latest

COPY ./ ./

RUN go build -o ./bin/SellerBot ./cmd

RUN mkdir "config"
RUN touch "config/config.yaml"

CMD ["./bin/SellerBot"]

