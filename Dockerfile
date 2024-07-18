FROM debian:12-slim

WORKDIR /app

COPY ./bin/qcg-center .
COPY ./assets/ ./assets/

ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

CMD ["bash", "-c", "./qcg-center"]