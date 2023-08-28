from golang:1.20 AS build

WORKDIR /usr/src/app

COPY go.mod . 
COPY go.sum .

RUN go mod download && \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app cmd/main.go

FROM golang:1.20-alpine

COPY --from=build /usr/src/app/bin/app .

CMD ["./app"]