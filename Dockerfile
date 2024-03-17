# First stage: Build the binary
FROM golang:latest AS build

# Set the working directory in the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.build.mod ./go.mod
COPY go.sum ./
COPY cmd/ ./cmd/
COPY internal ./internal
COPY .env .env
# Copy the source code to the container
# COPY . .
COPY scs/libsqlstore /go/src/scs/libsqlstore
# RUN ls -a
# RUN ls scs
# RUN ls scs/libsqlstore
# Download dependencies
RUN go mod vendor
# RUN ls vendor/github.com/mattn/go-sqlite3
# RUN cat vendor/github.com/mattn/go-sqlite3/error.go
# RUN go install ./cmd
RUN echo "APP_STAGE= STAGING" > app.env
# Build the binary
RUN go build -o main ./cmd

# Second stage: Copy the binary to a minimal image
# FROM alpine:latest
# RUN apk add --no-cache tzdata
# # Set the working directory
# WORKDIR /
# # Copy the binary from the build stage
# COPY --from=build /app/main .

# COPY --from=build /app/app.env .
# RUN ls 
# Make the binary executable
RUN chmod +x main 

# Expose the port that the application will run on
EXPOSE 8080

# Run the binary
ENTRYPOINT [ "./main" ]