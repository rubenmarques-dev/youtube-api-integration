# youtube-api-integration

A minimal Go HTTP service that wraps the YouTube Data API v3 and exposes a single endpoint to fetch video metadata.

Built following **Hexagonal Architecture** principles: the domain is free of external dependencies, adapters implement the domain ports, and wiring is done manually in `main.go`.

---

## Prerequisites

- Go 1.22+
- A valid [YouTube Data API v3](https://developers.google.com/youtube/v3/getting-started) key
- [Air](https://github.com/air-verse/air) for hot reload: `go install github.com/air-verse/air@latest`
- [swag CLI](https://github.com/swaggo/swag) for Swagger generation: `go install github.com/swaggo/swag/cmd/swag@latest`

---

## Running Locally

```bash
export YOUTUBE_API_KEY=your_api_key_here
export PORT=8080  # optional, defaults to 8080

go run ./cmd/server
```

---

## Running with Hot Reload (Air)

Air watches for file changes, regenerates Swagger docs, and restarts the server automatically.

```bash
export YOUTUBE_API_KEY=your_api_key_here
air
```

On each save, Air will:
1. Run `swag init` to regenerate Swagger docs
2. Recompile the server
3. Restart the process

---

## Swagger UI

Once the server is running, the interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

To manually regenerate the Swagger docs without running the server:

```bash
swag init -g cmd/server/main.go --output docs
```

---

## Running with Docker

```bash
docker build -t youtube-api-integration .

docker run -p 8080:8080 \
  -e YOUTUBE_API_KEY=your_api_key_here \
  youtube-api-integration
```

---

## Example Request

```bash
curl "http://localhost:8080/video?id=dQw4w9WgXcQ"
```

### Success Response (200)

```json
{
  "title": "Rick Astley - Never Gonna Give You Up (Official Music Video)",
  "description": "...",
  "viewCount": "1400000000",
  "likeCount": "15000000",
  "channelTitle": "Rick Astley",
  "publishedAt": "2009-10-25T06:57:33Z"
}
```

### Error Response (400)

```json
{ "error": "missing required query parameter: id" }
```

### Error Response (404)

```json
{ "error": "video not found" }
```
