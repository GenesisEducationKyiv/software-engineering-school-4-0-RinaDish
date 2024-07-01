FROM golang:1.22 as build

WORKDIR /app

COPY . .
RUN make build

WORKDIR /app/bin

CMD ["./main"]