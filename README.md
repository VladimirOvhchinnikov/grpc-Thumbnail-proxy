# gRPC Thumbnail Proxy & CLI Downloader

This project consists of a **gRPC proxy service** for fetching and caching YouTube video thumbnails and a **CLI utility** for requesting thumbnails via gRPC.

## Features

- gRPC proxy service to fetch and cache thumbnails
- CLI tool to request thumbnails from YouTube
- Supports SQLite for persistent caching
- CLI supports `--async` mode for parallel downloads

## Installation

### Prerequisites

- Go 1.19+
- SQLite 

### Clone the repository

```sh
git clone https://github.com/VladimirOvhchinnikov/grpc-thumbnail-proxy.git
cd grpc-thumbnail-proxy
```

## Building and Running the gRPC Server

Navigate to the service directory:

```sh
cd service
```

Build the binary:

```sh
go build -o grpc-thumbnail-server main.go
```

Run the server:

```sh
./grpc-thumbnail-server
```

By default, the server runs on port `50051`.

### Accessing Logs

Logs are written to `server.log` in the `service` directory by default. To view the logs, you can use the following commands:

- View the entire log file:

  ```sh
  cat server.log
  ```

- Follow the logs in real-time:

  ```sh
  tail -f server.log
  ```

If you want to change the log file location, update the configuration in `main.go`.

## Building and Using the CLI Tool

Navigate to the CLI directory:

```sh
cd cli
```

Build the binary:

```sh
go build -o grpc-thumbnail-cli main.go
```

Run the CLI to download a single thumbnail:

```sh
./grpc-thumbnail-cli -links "https://www.youtube.com/watch?v=EXAMPLE"
```

### Download multiple thumbnails

```sh
./grpc-thumbnail-cli -links "https://www.youtube.com/watch?v=EX1,https://www.youtube.com/watch?v=EX2"
```

### Download multiple thumbnails asynchronously

```sh
./grpc-thumbnail-cli -links "https://www.youtube.com/watch?v=EX1,https://www.youtube.com/watch?v=EX2 -async"
```

### CLI Help

To see available options, run:

```sh
./grpc-thumbnail-cli -h
```

## Running Tests

Tests are available for both the CLI and the service.

### Run all tests

```sh
go test ./...
```

### Specific tests

Run tests for the CLI:

```sh
cd cli
go test ./...

```

Run tests for the service:

```sh
cd service
go test ./...
```

## API Specification (gRPC)

Refer to `proto/service.proto` for details.

## Notes

- The repository does not contain compiled binaries. Build them using `go build`.
- Ensure the gRPC server is running before using the CLI.

## License

MIT Licens

