FROM alpine:edge

RUN apk add --no-cache chromium tor go
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o utr

CMD ["./utr", "--url", "https://www.youtube.com/watch?v=jQqnqA3KAjE", "--views", "1", "--instances", "1", "--duration", "10"]
