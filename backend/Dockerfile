FROM golang:1.21.5 

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .env ./

EXPOSE 2999

RUN CGO_ENABLED=0 GOOS=linux go build -o bin

CMD ["./bin", "--seed true"]


