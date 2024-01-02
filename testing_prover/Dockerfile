FROM golang:1.21-alpine as builder

WORKDIR /src
COPY . .

RUN apk --no-cache add ca-certificates

RUN go build -o /bin/app /src/main.go

EXPOSE 3000

ENTRYPOINT ["/bin/app"]
