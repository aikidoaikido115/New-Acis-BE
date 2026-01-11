FROM golang:1.24
WORKDIR /app

RUN curl -sSf https://atlasgo.sh | sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# รัน migration และ start server
CMD ["sh", "-c", "atlas migrate apply --env dev && ./main"]