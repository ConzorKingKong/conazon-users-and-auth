FROM golang:1.22.1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go .env ./
RUN go build .
EXPOSE 8080
CMD ["./conazon-users-and-auth"]