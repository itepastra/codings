FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY huffman ./huffman

RUN go build -o /codings

ENTRYPOINT ["/codings"]
