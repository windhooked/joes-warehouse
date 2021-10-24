FROM golang:1.17.2-alpine3.14

# Creating the `app` directory in which the app will run 
RUN mkdir /app

# Move everything from root to the newly created app directory
ADD . /app

# Specifying app as our work directory in which
# futher instructions should run into
WORKDIR /app

ENV CGO_ENABLED=0

# Download all neededed project dependencies
RUN go mod download

# Build the project executable binary
RUN go build -o main ./cmd/joes-warehouse

EXPOSE 7000/tcp

COPY ./cmd/joes-warehouse/config.toml /app/config.toml

# Run/Starts the app executable binary
CMD ["/app/main"]