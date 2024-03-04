# Mock HTTP Server

This repository provides a Golang implementation of a mock HTTP server that allows you to easily define and manage mock endpoints through a configuration YAML file. It also offers various command-line flags for additional configuration.

## Features

- **YAML Configuration:** Configure mock endpoints, including expected requests, responses, and more, using a user-friendly YAML file.
- **Multiple Endpoints:** Define multiple mock endpoints with different behaviors and responses.
- **Customizable Responses:** Set status codes, headers, and body content for each response to simulate various scenarios.
- **Command-Line Flags:** Override configuration options like port and TLS settings using command-line flags.
- **Basic Authentication:** Secure your mock server by enabling basic authentication with username and password flags.
- **TLS Support:** Enable secure communication over HTTPS by providing TLS certificate and key files.
- **Metrics:** Expose default go runtime metrics like request count, latency etc at `/metrics` endpoint

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/wasim-nihal/mock-http-server.git
    ```


2. Install dependencies:

    ```bash
    go mod download
    ```


## Usage

1. Build the executable:

    ```bash
    go build -o mock-server .
    ```


2. Run the server:

    ```bash
    ./mock-server \
      -server.config config.yaml \  # Path to your configuration file
      -port 8081 \                   # Optional port (default: 8080)
      -basicauth.username user      # Optional basic auth username
      -basicauth.password password  # Optional basic auth password
      -tls                          # Optional flag to enable TLS (requires certificate and key)
      -tlsCertFile cert.pem         # Optional path to TLS certificate file (if -tls is set)
      -tlsKeyFile key.pem           # Optional path to TLS key file (if -tls is set)
    ```


## Configuration File (YAML)

The YAML file defines the mock endpoints and their expected behavior. Here's an example structure:

```yaml
# List of mock endpoints
endpoints:
  # endpoint path
  /hello:
       # http method type
    -  method: GET
       # response content type
       content: text/plain
       # response body
       body: Hello, World! (GET)
       # response status
       status: 200

    -  method: POST
       content: text/plain
       body: Hello, World! (POST)
       status: 200
  /bye:
    -  method: GET
       content: text/plain
       body: Goodbye!
       status: 400
  # Add more endpoints with their respective configurations
  - ...
```

## License
This project is licensed under the Apache-2.0 license License. See the [LICENSE](https://github.com/wasim-nihal/mock-http-server/blob/master/LICENSE) file for details.

## Additional Notes
- Replace `config.yaml` with the actual name of your configuration file.
- Adjust the command-line flags and YAML example as necessary based on your specific implementation.
