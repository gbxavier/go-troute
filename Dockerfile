FROM golang:1.12-stretch

LABEL maintainer="Guilherme Bruno da Silva Xavier <gbxavier11@gmail.com>"

ARG APP_NAME=traceroute-gbxavier11

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/traceroute-gbxavier11

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Run the binary program produced by `go install`
ENTRYPOINT ["traceroute-gbxavier11"]