FROM golang:1.23

WORKDIR /app/backend


COPY go.mod go.sum ./

RUN go mod download && go mod verify 

COPY . .

RUN go build -v -o server

CMD ["/app/backend/server"]