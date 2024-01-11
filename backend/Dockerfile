FROM golang:1.21.5 

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .env Makefile ./

EXPOSE 2999

RUN go build -o bin

ENTRYPOINT ["/server/bin"]

CMD ["echo", "hello world"]
#CMD ["--seed", "true"]
