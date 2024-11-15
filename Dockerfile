# syntax=docker/dockerfile:1

FROM golang:1.23

# Add maintainer info
LABEL maintainer="jerryholland00@gmail.com"
LABEL version="0.1.0"

# create working directory
WORKDIR /jerry_holland_fetch_interview

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-jh-fetch-interview

EXPOSE 8080

# Run
CMD ["/docker-jh-fetch-interview"]