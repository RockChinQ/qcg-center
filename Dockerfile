FROM golang:1.21.7-bullseye

WORKDIR /app/qcg-center

COPY . .

RUN go mod download

RUN go build -o qcg-center src/main.go

RUN chmod +x ./qcg-center

ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

CMD ["./qcg-center"]