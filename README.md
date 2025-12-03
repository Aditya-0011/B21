# Proxy Server with TOTP Authentication

A secure HTTP proxy server written in Go that allows authenticated file downloads through a proxy endpoint. The server uses Time-based One-Time Password (TOTP) authentication to ensure only authorized users can access the proxy functionality.

## Features

- **TOTP Authentication**: Secure access control using time-based one-time passwords
- **Proxy Downloads**: Download files from any URL through the proxy server
- **Request Logging**: All requests and operations are logged with timestamps
- **Log Access**: Retrieve server logs through authenticated endpoint
- **Thread-Safe Logging**: Concurrent-safe logging with mutex protection

## Project Structure

```
.
├── cmd/
│   └── keygen.go          # TOTP key and QR code generator
├── server/
│   └── handlers.go        # HTTP request handlers
├── totp/
│   └── totp.go            # TOTP validation logic
├── utils/
│   └── logger.go          # Thread-safe logging utility
├── main.go                # Application entry point
├── go.mod                 # Go module dependencies
└── README.md              # Project documentation
```

## Prerequisites

- Go 1.25.3 or higher
- Environment variable `TOTP_SECRET` configured

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd proxy-server
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Configure the key generator:

   You can customize the TOTP metadata via environment variables (used by `cmd/keygen.go`):

   - `ISSUER`: Server/service name (default: `Server`)
   - `ACCOUNT_NAME`: Account identifier (default: `Admin`)

   Examples:

   - Windows PowerShell

     ```powershell
     $env:ISSUER="MyProxy"; $env:ACCOUNT_NAME="admin@example.com"; go run .\cmd\keygen.go
     ```

   - Linux/macOS

     ```bash
     ISSUER="MyProxy" ACCOUNT_NAME="admin@example.com" go run ./cmd/keygen.go
     ```

4. Generate TOTP secret and QR code:

   ```bash
   go run cmd/keygen.go
   ```

   This will generate:

   - A TOTP secret key (save this as an environment variable)
   - An ASCII QR code printed in the terminal (you can scan it directly)
   - A QR code image (`code.png`) that can be scanned with an authenticator app

   If you've already built the key generator binary, run it directly:

   - Linux/macOS: `./keygen`
   - Windows: `.\keygen.exe`

5. Set up your authenticator app:

   You can add the TOTP entry to your authenticator app (Google Authenticator, Authy, etc.) in two ways:

   - **Scan QR Code**: Scan the QR printed in the terminal, or open the generated `code.png` file and scan it with your authenticator app
   - **Manual Entry**: Alternatively, manually enter the secret key from step 4 into your authenticator app along with the issuer and account name

6. Set the environment variable:

   ```bash
   # Windows PowerShell
   $env:TOTP_SECRET="your-generated-secret"

   # Windows CMD
   set TOTP_SECRET=your-generated-secret

   # Linux/macOS
   export TOTP_SECRET=your-generated-secret
   ```

## Usage

### Build

- Server binary:

  - Bash/macOS/Linux:

    ```bash
    GOOS=linux GOARCH=arm64 go build -o proxy main.go
    ```

  - Windows PowerShell:

    ```powershell
    $env:GOOS="linux"; $env:GOARCH="arm64"; go build -o proxy main.go
    ```

- Key generator binary:

  - Bash/macOS/Linux:

    ```bash
    GOOS=linux GOARCH=arm64 go build -o keygen ./cmd/keygen.go
    ```

  - Windows PowerShell:

    ```powershell
    $env:GOOS="linux"; $env:GOARCH="arm64"; go build -o keygen ./cmd/keygen.go
    ```

Note: Adjust `GOOS`/`GOARCH` for your target platform as needed.

### Starting the Server

Run the server:

```bash
go run main.go
```

The server will start on port `7000` by default (configurable via `PORT`).

If you've built the server binary, run it directly:

- Linux/macOS: `./proxy`
- Windows: `.\proxy.exe`

### API Endpoints

#### 1. Download Files (Proxy Endpoint)

**Endpoint**: `POST /`

**Description**: Download a file from a specified URL through the proxy server.

**Request Body**:

```json
{
  "url": "https://example.com/file.zip",
  "otp": "123456"
}
```

**Response**: The requested file content with appropriate headers.

**Example**:

```bash
curl -X POST http://localhost:6000/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/file.zip","otp":"123456"}'
```

#### 2. Retrieve Server Logs

**Endpoint**: `POST /logs`

**Description**: Download the server log file (requires TOTP authentication).

**Request Body**:

```json
{
  "otp": "123456"
}
```

**Response**: Server log file (`proxy.log`)

**Example**:

```bash
curl -X POST http://localhost:6000/logs \
  -H "Content-Type: application/json" \
  -d '{"otp":"123456"}' \
  --output logs.txt
```

## Security

- **TOTP Authentication**: All endpoints require a valid TOTP code generated from your authenticator app
- **Request Logging**: All authentication attempts and operations are logged
- **Thread-Safe Operations**: Log file access is protected with mutexes to prevent corruption

## Configuration

Environment variables used by the server:

- `TOTP_SECRET`: Required TOTP secret for validating OTPs
- `LOG_FILE`: Log file path (default: `proxy.log`)
- `PORT`: Server port (default: `7000`)
- `ISSUER`: Server/service name (default: `Server`)
- `ACCOUNT_NAME`: Account identifier (default: `Admin`)

## Error Handling

The server handles various error scenarios:

- **Invalid Method**: Only POST requests are accepted
- **Bad Request**: Malformed JSON in request body
- **Invalid OTP**: TOTP validation failure
- **Server Error**: Missing TOTP secret configuration
- **Bad Gateway**: Unable to reach target URL
- **Internal Server Error**: File system or logging errors

## Logging

All operations are logged to both:

- Console output
- `proxy.log` file

By default, `proxy.log` is created in the process's current working directory.

Log entries include:

- Server startup
- Authentication attempts (success/failure)
- Download requests with IP addresses
- Data transfer statistics
- Errors and exceptions

## Dependencies

- `github.com/pquerna/otp` - TOTP generation and validation
- `github.com/boombuler/barcode` - QR code generation for TOTP setup
- `github.com/mdp/qrterminal/v3` - Terminal QR code rendering

## Author

[Aditya](https://adityapunmiya.com)

## Deployment

Deployment and CI/CD are environment-specific. Use the build commands above within your own pipeline or tooling appropriate for your infrastructure.
