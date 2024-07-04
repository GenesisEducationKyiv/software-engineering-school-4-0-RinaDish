FROM golang:1.22

WORKDIR /app

COPY . .
RUN make build

WORKDIR /app/bin

CMD ["./main"]