### WRAITH SMS DOCKERFILE ###

# Start with the official Alpine base image
FROM alpine:latest

# Set build-time args and env
ARG workdir=/wraith_sms

# Install necessary packages
RUN apk add --no-cache go make

# Create work directories
RUN mkdir -p $workdir

# Copy the repo and configs into the container
COPY ./message_server $workdir
RUN cd $workdir && \
rm -f secrets.env config.toml

# Compile the project
RUN cd $workdir && \
make build

# Expose necessary ports
EXPOSE 8080

# Run the binary
WORKDIR $workdir
CMD ["./sms_server"]