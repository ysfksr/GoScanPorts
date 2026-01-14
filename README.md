# GoScanPorts

A fast and efficient port scanner written in Go with concurrent scanning capabilities.

## Features

- Concurrent port scanning for fast results
- Customizable port range
- Configurable connection timeout
- Common service name detection
- Clean console output

## Installation

```bash
go build -o portscan main.go
```

## Usage

Basic usage:
```bash
go run main.go -host localhost -start 1 -end 1024
```

### Command Line Options

- `-host` - Host to scan (default: localhost)
- `-start` - Start port (default: 1)
- `-end` - End port (default: 1024)
- `-timeout` - Connection timeout in milliseconds (default: 1000)
- `-popular` - Scan only popular ports for web, databases, remote access, etc. (default: false)

### Examples

Scan localhost ports 1-1000:
```bash
go run main.go -host localhost -start 1 -end 1000
```

Scan a specific IP address:
```bash
go run main.go -host 192.168.1.1 -start 1 -end 65535
```

Quick scan of common ports:
```bash
go 

**Scan only popular ports** (recommended for quick checks):
```bash
go run main.go -host localhost -popular
```

This scans ~40 common ports including:
- Web servers (80, 443, 8080, 8443, 3000, 4200, 5000, 8000)
- Databases (MySQL, PostgreSQL, MongoDB, Redis, MS SQL, CouchDB, Elasticsearch)
- Remote access (SSH, FTP, RDP, VNC)
- Email services (SMTP, POP3, IMAP and their secure variants)
- Other services (DNS, SMB, RabbitMQ, Prometheus, Memcached)run main.go -host example.com -start 1 -end 1024 -timeout 500
```

## Sample Output

```
Scanning localhost from port 1 to 1024...
This may take a while depending on the range...

Scan completed in 2.3s
Found 3 open port(s):
  Port 22 is open (SSH)
  Port 80 is open (HTTP)
  Port 443 is open (HTTPS)
```

## Note

Use this tool responsibly. Only scan hosts you own or have permission to scan.