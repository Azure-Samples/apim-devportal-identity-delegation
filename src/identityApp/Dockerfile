FROM --platform=linux/amd64 golang:1.19.3-alpine3.17 

# Define current working directory
WORKDIR /auth-proxy

# Download modules to local cache so we can skip re-
# downloading on consecutive docker build commands
COPY src/identityApp/go.mod .
COPY src/identityApp/go.sum .
RUN go mod download

# Add sources
COPY src/identityApp .
COPY .env* /auth-proxy/

RUN go build -o out/auth0-go-web-app .

# Expose port 3000 for our web app binary
EXPOSE 3000

CMD ["/auth-proxy/out/auth0-go-web-app"]
