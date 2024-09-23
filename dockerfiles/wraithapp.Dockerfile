### WRAITH SMS DOCKERFILE ###

# Start with the official Alpine base image
FROM alpine:latest

# Set build-time args and env
ARG workdir=/wraith_sms

# Install necessary packages
RUN apk add --no-cache go make

# Set GOPATH and add it to PATH
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Create work directories
RUN mkdir -p $workdir

# Set the working directory
WORKDIR $workdir

# Copy only the go.mod and go.sum files (if they exist)
COPY ./message_server/go.mod ./message_server/go.sum* ./

# Download dependencies
RUN go mod download

# Expose necessary ports
EXPOSE 8080

# Install air for hot reloading; using an older version that's compatible with Alpine's packaged version
RUN go install github.com/cosmtrek/air@v1.49.0

# Verify air installation
RUN which air

# Copy the .air.toml file into the container
COPY ./message_server/.air.toml .

# Set the entrypoint to use air for hot reloading
ENTRYPOINT ["air", "-c", ".air.toml"]

