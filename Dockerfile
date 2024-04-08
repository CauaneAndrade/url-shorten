FROM golang:1.22
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["./main"]
