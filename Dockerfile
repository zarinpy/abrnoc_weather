FROM golang:1.25.2-alpine
WORKDIR /app
ENV GOTOOLCHAIN=auto
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["./main"]
