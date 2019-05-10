FROM golang:alpine

WORKDIR /Users/nanda/fogfing-repository/material/greet/greet_server

COPY . .

RUN apk add --no--cache bash git openssh

RUN go -d -v ./...

RUN go get ./...

EXPOSE 50051

CMD ["go", "run", "server.go"]