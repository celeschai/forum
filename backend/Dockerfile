FROM golang:1.21.5 

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY *.go Makefile ./

EXPOSE 2000

RUN go build -o bin

CMD ["./bin", "--seed", "true"]
