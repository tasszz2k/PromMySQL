# Use the official Golang image as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Golang application
RUN go build -o prommysql .

# Expose the port that Prometheus will scrape
EXPOSE 8080

# Set the entry point for the container
ENTRYPOINT ["./prommysql"]
