# Build stage.
FROM golang:1.17-alpine as builder

# Add git tool to extract latest commit and tag.
# Add root certificates to be used for ssl/tls.
# Add openssl to build self-signed certificates.
RUN apk add --update --no-cache ca-certificates git openssl

# Setup the working directory
WORKDIR /app/

# Copy go mod file and download dependencies.
COPY go.* ./
# RUN go mod download -x

# Copy all files to the container’s workspace.
COPY . .

# Execute the self-signed certificate generation script.
RUN chmod +x ./scripts/generate.certs.sh
RUN ./scripts/generate.certs.sh

# Build the app program inside the container.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o demo-rest-api-server -a -ldflags "-extldflags '-static' -X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" .

# Final stage with minimalist image.
FROM scratch

LABEL maintainer="Jerome Amon"

# Copy our static executable to the new container root.
COPY --from=builder ./app/demo-rest-api-server ./demo-rest-api-server

# Copy self-signed certs and ca-certs and configuration and databases migrations files.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder ./app/server.crt ./assets/certs/server.crt
COPY --from=builder ./app/server.key ./assets/certs/server.key
COPY --from=builder ./app/server.config.docker.yml ./server.config.yml
COPY --from=builder ./app/pkg/infrastructure/postgres/migrations ./pkg/infrastructure/postgres/migrations
COPY --from=builder ./app/pkg/infrastructure/mongo/migrations ./pkg/infrastructure/mongo/migrations

EXPOSE 8080
ENTRYPOINT [ "./demo-rest-api-server" ]
