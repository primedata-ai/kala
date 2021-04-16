FROM golang

WORKDIR /go/src/github.com/primedata-ai/kala
COPY . .
RUN go build && mv kala /usr/bin

CMD ["kala", "serve"]
EXPOSE 8000
