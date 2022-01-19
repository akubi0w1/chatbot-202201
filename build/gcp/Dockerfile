FROM golang:1.16

ADD bin/gcp/chatbot /go/bin

RUN chmod 777 /go/bin/chatbot

ENV TZ "Asia/Tokyo"

CMD ["./bin/chatbot"]
EXPOSE 8080