FROM golang:1.18-alpine

RUN apk -u add make

WORKDIR /usr/src/github.com/downloop/api

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make BINARY=/usr/local/bin/api

EXPOSE 8080
CMD ["api"]