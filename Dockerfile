FROM golang:1.18.10-bullseye

WORKDIR /app/qcg-center

COPY . .

RUN go mod download

RUN go build -o qcg-center src/main.go

RUN chmod +x ./qcg-center

CMD ["./qcg-center"]