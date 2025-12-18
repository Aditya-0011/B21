# ✈️ Project B21 (Stealth Relay)

**A lightweight, authenticated artifact relay service designed for high-security environments.**

---

## 📖 About

**Project B21** is a secure, authenticated streaming relay written in Go. It enables the retrieval of external artifacts (updates, binaries, patches) into restrictive network environments where direct access to third-party domains is blocked or unstable.

Unlike standard proxies, B21 does not expose a full internet gateway. Instead, it creates a **single-purpose, verified tunnel** for authorized file transfers using strict "Defense in Depth" mechanisms. Its primary design goal is **stealth** and **minimal footprint**.

---

## ✨ Key Features

- 🔐 **TOTP Authentication** - Zero-trust access control using Time-based One-Time Passwords (RFC 6238)
- ⚡ **Streaming Proxy** - Downloads files from any URL using efficient `io.Reader/Writer` pipes (constant low RAM usage)
- 🕵️ **Stealth Operations** - Masks traffic as generic binary data (`application/octet-stream`) to bypass file-type filters
- ⏯️ **Resumable Downloads** - Full support for `Range` headers, allowing large file transfers to be paused, resumed, or chunked
- 📝 **Audit Logging** - Comprehensive request logging with timestamping and Real IP detection (via `X-Forwarded-For`)
- 🛡️ **Thread-Safe** - Concurrent-safe logging with mutex protection for reliability under load
- 📦 **Minimal Dependencies** - Compiles into a single, statically-linked binary
- 🌐 **Connection Testing** - Simple GET endpoint to verify server availability

---

## 📁 Project Structure

```
.
├── .github/
│   └── workflows/
│       ├── keygen.yml     # CI/CD for key generator deployment
│       └── server.yml     # CI/CD for server deployment
├── cmd/
│   └── keygen.go          # Utility to generate TOTP secrets and QR codes
├── server/
│   └── handlers.go        # Core HTTP request logic (Proxy, Logs, Range Support)
├── totp/
│   └── totp.go            # TOTP validation wrapper
├── utils/
│   └── logger.go          # Thread-safe logging utility
├── .gitattributes         # Git configuration for line endings
├── .gitignore             # Git ignore patterns
├── go.mod                 # Go module dependencies
├── main.go                # Application entry point
└── README.md              # Documentation (this file)
```

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.25.3** or higher (required only for development and building binaries)
- An authenticator app (Google Authenticator, Authy, 1Password, etc.)
- Access to a VPS or external server to host the relay (optional for production)

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Aditya-0011/B21
   cd B21
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

---

## 🔑 Configuration & Setup

Before running the server, you must generate a secure TOTP secret.

### Step 1: Generate TOTP Secret

The [`keygen.go`](cmd/keygen.go) utility generates a TOTP secret and QR code for use with authenticator apps.

**Linux / macOS:**

```bash
ISSUER="B21-Relay" ACCOUNT_NAME="Admin" go run cmd/keygen.go
```

**Windows (PowerShell):**

```powershell
$env:ISSUER="B21-Relay"; $env:ACCOUNT_NAME="Admin"; go run .\cmd\keygen.go
```

**Output:**

- A **Secret Key** (save this securely!)
- A **QR Code** printed in the terminal (scan with your authenticator app)
- A `code.png` file (backup image)

### Step 2: Set Environment Variables

Configure the server environment before running.

| Variable       | Description                           | Default     | Required |
| :------------- | :------------------------------------ | :---------- | :------- |
| `TOTP_SECRET`  | The secret generated in Step 1        | N/A         | ✅ Yes   |
| `PORT`         | The port the server listens on        | `7000`      | ❌ No    |
| `LOG_FILE`     | Path to the log file                  | `proxy.log` | ❌ No    |
| `ISSUER`       | Server/service name (for keygen only) | `Server`    | ❌ No    |
| `ACCOUNT_NAME` | Account identifier (for keygen only)  | `Admin`     | ❌ No    |

**Example:**

**Linux / macOS:**

```bash
export TOTP_SECRET="YOUR_SECRET_HERE"
export PORT="8080"
export LOG_FILE="/var/log/b21.log"
```

**Windows (PowerShell):**

```powershell
$env:TOTP_SECRET="YOUR_SECRET_HERE"
$env:PORT="8080"
$env:LOG_FILE="C:\logs\b21.log"
```

---

## 🛠️ Usage

### Running the Server

#### Development Mode

Run directly with Go:

```bash
go run main.go
```

#### Production Build

For stealth deployments, strip debug symbols using `-ldflags "-s -w"` to reduce binary size.

**Linux / macOS:**

```bash
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o b21 main.go
```

**Windows** (if development environment is Windows):

```powershell
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -ldflags "-s -w" -o b21 main.go
```

> 💡 **Tip:** Adjust `GOOS` and `GOARCH` for your target platform. Common combinations:
>
> - Linux AMD64: `GOOS=linux GOARCH=amd64`
> - Linux ARM64: `GOOS=linux GOARCH=arm64`
> - Windows AMD64: `GOOS=windows GOARCH=amd64`
> - macOS ARM64: `GOOS=darwin GOARCH=arm64`

---

## 📡 API Endpoints

### 1. 🏥 Health Check / Connection Test

Simple endpoint to verify server availability.

- **Method:** `GET`
- **URL:** `/`
- **Authentication:** None required
- **Response:** Current server timestamp in RFC3339 format

**Example:**

```bash
curl https://example.com/
```

**Response:**

```
2024-01-15T10:30:45Z
```

---

### 2. 📥 Download File (Proxy Endpoint)

Relays a file from a remote URL to the client. Supports `Range` headers for resumable downloads.

- **Method:** `POST`
- **URL:** `/`
- **Headers:** `Content-Type: application/json`
- **Authentication:** TOTP required

**Stealth Behavior:** The server forces the response filename to `resource.dat` and content-type to `application/octet-stream` to avoid firewall blocking. You must rename the file locally.

**Request Body:**

```json
{
  "url": "https://example-cdn.com/file.zip",
  "otp": "123456"
}
```

**Example:**

```bash
# Basic download
curl -X POST https://example.com/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example-cdn.com/file.zip","otp":"123456"}' \
  -o downloaded-file.zip

# Download with rate limiting (recommended for stealth)
curl --limit-rate 2M -X POST https://example.com/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example-cdn.com/file.zip","otp":"123456"}' \
  -o large-file.zip
```

**cURL Options Explained:**

- `--limit-rate 2M` - Limit bandwidth to 2MB/s (blends with video traffic patterns)
- `-o filename` - Save output to specified file (rename from `resource.dat`)

#### 🔄 Resuming Interrupted Downloads

**Important:** curl's `-C -` flag does **not** work with POST requests. To resume an interrupted download, you must manually specify the byte offset.

**Quick Resume (Practical Example):**

Let's say your `some-secret.zip` download got interrupted. Here's how to resume:

**Windows (PowerShell):**

```powershell
# Step 1: Check your partial file size
$bytes = (Get-Item some-secret.zip).Length
Write-Host "File size: $bytes bytes - Use this in Range header"

# Step 2: Resume download (replace YOUR_BYTE_SIZE with the number from above)
curl --limit-rate 2M -X POST https://example.com/ `
  -H "Content-Type: application/json" `
  -H "Range: bytes=$bytes-" `
  -d "{`"url`":`"https://example-cdn.com/large-file.zip`",`"otp`":`"123456`"}" `
  --output some-secret.zip
```

**Linux/macOS (Bash):**

```bash
# Step 1: Check your partial file size
bytes=$(stat -c%s some-secret.zip)  # Linux
# bytes=$(stat -f%z some-secret.zip)  # macOS
echo "File size: $bytes bytes - Use this in Range header"

# Step 2: Resume download (the $bytes variable is automatically used)
curl --limit-rate 2M -X POST https://example.com/ \
  -H "Content-Type: application/json" \
  -H "Range: bytes=$bytes-" \
  -d '{"url":"https://example-cdn.com/large-file.zip","otp":"123456"}' \
  -o some-secret.zip
```

**Wait, what about overwriting?**

⚠️ **Important:** The `-o` (output) flag in curl will **overwrite** your partial file if you don't handle it carefully. Here are two safe approaches:

**Option 1: Use a different filename, then combine**

**Windows (PowerShell):**

```powershell
# Download remaining part to a new file
curl ... -o some-secret-part2.zip

# Then combine files
Get-Content some-secret.zip, some-secret-part2.zip -Encoding Byte | Set-Content some-secret-complete.zip -Encoding Byte
```

**Linux/macOS (Bash):**

```bash
# Download remaining part to a new file
curl ... -o some-secret-part2.zip

# Then combine files
cat some-secret.zip some-secret-part2.zip > some-secret-complete.zip
```

**Option 2: Let curl write directly (but backup first)**

**Windows (PowerShell):**

```powershell
# Make a backup just in case
Copy-Item some-secret.zip some-secret-backup.zip

# Now resume - curl will handle the Range request properly
curl --limit-rate 2M -X POST https://example.com/ `
  -H "Range: bytes=123838464-" `
  ... `
  -o some-secret.zip
```

**Linux/macOS (Bash):**

```bash
# Make a backup just in case
cp some-secret.zip some-secret-backup.zip

# Now resume - curl will handle the Range request properly
curl --limit-rate 2M -X POST https://example.com/ \
  -H "Range: bytes=123838464-" \
  ... \
  -o some-secret.zip
```

> 💡 **Pro Tip:** When the server properly supports Range requests (which B21 does), and you specify the exact byte offset that matches your file size, curl should only receive the remaining bytes. However, always keep a backup or use separate files to be safe!

**How It Works:**

1. You check your partial file size (e.g., 50MB = 52428800 bytes)
2. You send a `Range: bytes=52428800-` header, telling the server to start from byte 52428800
3. B21 forwards this Range header to the upstream server
4. The upstream server returns only the remaining data with `206 Partial Content`
5. You append this data to your existing file, completing the download

> **Note:** Resume functionality requires the upstream URL to support HTTP Range requests. Most modern CDNs and file servers do.

**Response Headers:**

```
Content-Type: application/octet-stream
Content-Disposition: attachment; filename=resource.dat
Accept-Ranges: bytes
Content-Length: [file size]
```

---

### 3. 📋 Retrieve Server Logs

Securely fetches the server's internal activity logs via a non-blocking snapshot.

- **Method:** `POST`
- **URL:** `/logs`
- **Headers:** `Content-Type: application/json`
- **Authentication:** TOTP required

**Request Body:**

```json
{
  "otp": "123456"
}
```

**Example:**

```bash
curl -X POST https://example.com/logs \
  -H "Content-Type: application/json" \
  -d '{"otp":"123456"}' \
  -o server-logs.txt
```

**Log Contents:**

The log file contains timestamped entries including:

- Server startup events
- Authentication attempts (success/failure)
- Download requests with IP addresses
- Data transfer statistics
- Errors and exceptions

---

## 🛡️ Security Model

### Defense in Depth

1. **Authentication Layer** - Every request is validated against the `TOTP_SECRET`. No public access is allowed.
2. **Traffic Masking** - All outgoing files are masqueraded as generic binary data (`resource.dat` / `application/octet-stream`).
3. **Concurrency Protection** - Log files are protected by [`utils.LogMutex`](utils/logger.go) to prevent race conditions.
4. **Minimal Attack Surface** - Only POST/GET requests are accepted. No query parameters or complex headers are parsed.
5. **IP Logging** - Real client IPs are extracted via `X-Forwarded-For` header (supports reverse proxy scenarios).

---

## ⚠️ Error Handling

The server handles various error scenarios with appropriate HTTP status codes:

| Scenario                   | HTTP Status | Response                 | Logged As          |
| :------------------------- | :---------- | :----------------------- | :----------------- |
| Invalid HTTP method        | `405`       | "Method not allowed"     | `[INVALID METHOD]` |
| Invalid or expired OTP     | `403`       | "Invalid OTP"            | `[AUTH FAIL]`      |
| Missing TOTP secret        | `500`       | "Server error"           | `[ERROR]`          |
| Unable to reach target URL | `502`       | "Failed to reach target" | `[ERROR]`          |

---

## 📊 Logging

All operations are logged to **both** console output and the configured log file (default: [`proxy.log`](proxy.log)).

### Log Entry Types

- `[INFO]` - Informational messages (connection tests)
- `[PROXY]` - Download operations
- `[SUCCESS]` - Successful file transfers
- `[AUTH FAIL]` - Failed authentication attempts
- `[INVALID METHOD]` - Invalid HTTP methods
- `[ERROR]` - Server errors
- `[ADMIN]` - Log export operations

### Log Format

```
2024/01/15 10:30:45 [PROXY] Starting download: https://example.com/file.zip (IP: 192.168.1.100)
2024/01/15 10:31:02 [SUCCESS] Transferred 104857600 bytes
```

### Thread Safety

The [`utils.Logger`](utils/logger.go) uses a mutex ([`utils.LogMutex`](utils/logger.go)) to ensure thread-safe writes during concurrent requests.

---

## 📦 Dependencies

This project uses minimal external dependencies:

- [`github.com/pquerna/otp`](https://github.com/pquerna/otp) - TOTP generation and validation
- [`github.com/boombuler/barcode`](https://github.com/boombuler/barcode) - QR code generation for TOTP setup (keygen only)
- [`github.com/mdp/qrterminal/v3`](https://github.com/mdp/qrterminal) - Terminal QR code rendering (keygen only)

All dependencies are declared in [go.mod](go.mod).

---

## 🚀 Deployment

### CI/CD Workflows

This project includes GitHub Actions workflows for automated deployment:

- [`.github/workflows/server.yml`](.github/workflows/server.yml) - Builds and deploys the main server binary (`b21`)
- [`.github/workflows/keygen.yml`](.github/workflows/keygen.yml) - Builds and deploys the key generator utility

**Triggers:**

- Push tags matching `server*` for server deployment
- Push tags matching `keygen*` for keygen deployment

**Example:**

```bash
# Deploy server
git tag server-v1.0.0
git push origin server-v1.0.0

# Deploy keygen
git tag keygen-v1.0.0
git push origin keygen-v1.0.0
```

---

## 🗺️ Roadmap

- [ ] **Bandwidth Throttling** - Server-side rate limiting to enforce "Low and Slow" traffic profiles automatically
- [ ] **Client CLI** - A dedicated Go CLI tool to handle automated downloads, local file renaming, and seamless resume functionality (no manual Range headers needed)
- [ ] **Metrics Dashboard** - Real-time monitoring of transfer statistics

---

## 👤 Author

- [**Aditya Punmiya**](https://adityapunmiya.com)

---

**⚠️ Disclaimer:** Have fun with it.
