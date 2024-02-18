FROM golang:latest

COPY ./ ./

RUN go build -o ./bin/SellerBot ./cmd

CMD ["./bin/SellerBot"]
#FROM alpine:latest
#
#RUN apk --no-cache add ca-certificates
#
#COPY --from=0 ./bin/ ./bin/
#RUN chmod +x ./bin/
#
#CMD ["./bin/SellerBot"]

