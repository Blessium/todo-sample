FROM golang:buster

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o todo-sample

EXPOSE 1234

CMD ["/app/todo-sample"]
