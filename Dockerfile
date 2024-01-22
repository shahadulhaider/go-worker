FROM golang:1.20

WORKDIR /go/src/app

COPY . .

RUN go build -o app

ARG MONGO_URI
ARG DB_NAME
ARG COLLECTION_NAME
ARG CRON_INTERVAL

ENV MONGO_URI=$MONGO_URI
ENV DB_NAME=$DB_NAME
ENV COLLECTION_NAME=$COLLECTION_NAME
ENV CRON_INTERVAL=$CRON_INTERVAL

CMD ["./app"]