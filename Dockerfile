FROM golang:1.22.3

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev

CMD ["tail", "-f", "/dev/null"]
