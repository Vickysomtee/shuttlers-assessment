# Build Stage
    FROM golang:1.24.3 AS builder

    WORKDIR /app
    
    # Copy go files
    COPY go.mod ./
    RUN go mod download
    
    COPY . .
    
    RUN go build -o server .
    
    # Run Stage
    FROM gcr.io/distroless/base-debian12
    
    WORKDIR /app
    
    COPY --from=builder /app/server .

    EXPOSE 9877
    
    CMD ["./server"]
    