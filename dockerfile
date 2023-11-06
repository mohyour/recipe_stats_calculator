FROM golang:1.20-alpine
WORKDIR /src

COPY . .
RUN go mod download

RUN go build -o /recipe

EXPOSE 8080

ENTRYPOINT [ "/recipe" ]