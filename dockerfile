FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal/ ./internal
COPY ./cloud-storage-credential.json ./
RUN go build -o /go-template ./cmd

EXPOSE 8080

CMD ["sh", "-c", "/go-template"]