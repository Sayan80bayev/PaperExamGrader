FROM gocv/opencv:latest

WORKDIR /app

# Copy go.mod and go.sum from the root of your project
COPY go.mod go.sum ./

RUN go mod download

# Copy the entire project (including cmd/server)
COPY . .

# Build the server binary from the cmd/server directory
RUN go build -o server ./cmd/server

EXPOSE 8081

CMD ["./server"]