FROM golang:1.22-alpine

RUN apk add --no-cache sqlite
RUN apk add --no-cache gcc
RUN apk add --no-cache musl-dev
# RUN apk add --no-cache git

# RUN git clone https://github.com/luiz734/chatapp-api


WORKDIR chatapp-api
COPY . .
RUN rm go.sum go.mod
RUN go mod init chatapp-api
RUN go mod tidy

# CMD sqlite3
CMD go run .
