# youtube-api-integration

A minimal Go HTTP service that wraps the YouTube Data API v3 and exposes a single endpoint to fetch video metadata.

Built following **Hexagonal Architecture** principles: the domain is free of external dependencies, adapters implement the domain ports, and wiring is done manually in `main.go`.

---

## Prerequisites

- Go 1.22+
- A valid [YouTube Data API v3](https://developers.google.com/youtube/v3/getting-started) key

---

## Running Locally

```bash
export YOUTUBE_API_KEY=your_api_key_here
export PORT=8080  # optional, defaults to 8080

go run ./cmd/server
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
