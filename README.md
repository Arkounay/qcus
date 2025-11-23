# QCUS: Quick CLI Upload Server

A temporary cli-first file upload server

## Features

- **CLI upload support with cURL**
- **Easily self-hostable with Docker support** 
- ASCII qr code
- Password-protected uploads
- Real-time download notifications (WebSocket)
- One-click copy for URLs and cURL commands
- Auto-delete files after download or configurable expiry time
- Configurable via environment variables

## Demo

- https://qcus.outerark.com/
- password: `demo`

## Quick Start with Docker

```bash
docker run -d --restart unless-stopped -p 8088:8088 --name qcus arkounay/qcus
```

### Easy self-host example with Caddy

```bash
docker run -d --restart unless-stopped -p 127.0.0.1:8091:8088 -e UPLOAD_PASSWORD=mysecret --name qcus arkounay/qcus
```

**Caddyfile**

```caddyfile
example.com {
    reverse_proxy :8091
}
```

### With Custom Configuration

```bash
docker run -p 8088:8088 \
  -e UPLOAD_PASSWORD=mysecret \
  -e MAX_FILE_SIZE_MB=500 \
  -e FILE_EXPIRY_MINUTES=30 \
  arkounay/qcus
```

### With Persistent Storage (Optional)

```bash
docker run -p 8088:8088 \
  -v $(pwd)/uploads:/app/uploads \
  arkounay/qcus
```

## Setup (Local Development)

### 1. Install Node.js dependencies

```bash
npm install
```

### 2. Build the frontend for the web page

```bash
npm run build
```

### 3. Run the Go server

```bash
go run main.go
```

or for hot reload:

```bash
npm run dev
```

The Vite dev server will run on `http://localhost:5173` with hot module reloading.

### 4. Open in browser

Navigate to `http://localhost:8088`

Default password: `demo`


## Configuration

All configuration is done via environment variables:

| Variable              | Default | Description                      |
|-----------------------|---------|----------------------------------|
| `UPLOAD_PASSWORD`     | `demo`  | Password required for uploads    |
| `MAX_FILE_SIZE_MB`    | `100`   | Maximum file size in megabytes   |
| `FILE_EXPIRY_MINUTES` | `10`    | Minutes before files auto-delete |
| `PORT`                | `8088`  | Server port                      |

## CLI Usage

### Upload a file

```bash
curl -F "file=@yourfile.txt" -H "X-Upload-Password: demo" http://localhost:8088
```

### Download a file

```bash
curl -o "filename.txt" http://localhost:8088/download/{fileID}
```


## cli upload example

```bash
antonin@Antonin ~> curl -T test.txt -H "X-Upload-Password: demo" https://qcus.outerark.com
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      
                                                                      

File uploaded successfully!
Original name: test.txt
File size: 14 B
Download URL: http://qcus.outerark.com/download/a7496105fae5e95cef51aec0bf4f1a02
cURL command: curl -o "test.txt" http://qcus.outerark.com/download/a7496105fae5e95cef51aec0bf4f1a02

```