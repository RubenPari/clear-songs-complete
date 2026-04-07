# Clear Songs

Clear Songs is a powerful REST API service built with Go that helps you efficiently manage your Spotify music library by providing comprehensive bulk deletion capabilities for your tracks and playlists.

![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![API Version](https://img.shields.io/badge/API-v1.0-orange.svg)

## 🎯 Overview

The application allows you to quickly clean up your Spotify library by:

- Removing all tracks from a specific artist across your entire library
- Deleting tracks based on quantitative criteria (number of songs per artist)
- Clearing out entire playlists while maintaining the playlist structure
- Converting albums to individual tracks in your library
- Comprehensive backup system to prevent accidental data loss

## ✨ Features

### 🎵 Track Management

- **Delete by Artist**: Remove all tracks from a specific artist in your library
- **Quantitative Deletion**: Delete tracks based on the number of songs you have per artist (e.g., remove all tracks from artists with more than X songs)
- **Range-based Deletion**: Filter and delete tracks within specific count ranges
- **Track Analysis**: Get detailed summaries of your library organized by artist

### 📋 Playlist Management

- **Playlist Clearing**: Empty any playlist you own while keeping the playlist itself intact
- **Dual Deletion**: Remove tracks from both playlist and your personal library simultaneously
- **Bulk Operations**: Perform operations quickly and efficiently through the API

### 🔄 Album Operations

- **Album Conversion**: Convert albums to individual songs in your library

### 🛡️ Safety Features

- **Automatic Backup**: All deleted tracks are automatically saved to PostgreSQL database
- **Recovery System**: Restore accidentally deleted tracks from the backup
- **Smart Caching**: Intelligent cache management for optimal performance
- **Transaction Safety**: Operations are performed safely with proper error handling

## 🏗️ Technical Stack

- **Backend**: Go (Golang) 1.23+
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Caching**: In-memory cache with automatic invalidation
- **Authentication**: OAuth 2.0 with Spotify
- **Documentation**: Swagger/OpenAPI
- **External API**: Spotify Web API

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- Spotify Developer Account
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/RubenPari/clear-songs.git
   cd clear-songs
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Configure environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Set up your Spotify App**
   - Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
   - Create a new app
   - Add your redirect URI (e.g., `http://127.0.0.1:3000/auth/callback`)
   - Copy Client ID and Client Secret to your `.env` file

5. **Run the application**

   ```bash
   go run src/main.go
   ```

The server will start on `http://127.0.0.1:3000`

## ⚙️ Environment Configuration

Create a `.env` file in the root directory with the following parameters:

```env
# Spotify API Credentials
CLIENT_ID=your_spotify_client_id
CLIENT_SECRET=your_spotify_client_secret
REDIRECT_URL=http://127.0.0.1:3000/auth/callback

# Database Configuration
DB_HOST=127.0.0.1
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=clear_songs
DB_PORT=5432

# Redis Cache Configuration
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 🐳 Docker Setup

### Using Docker Compose (Recommended)

The easiest way to run the entire backend stack (API, PostgreSQL, and Redis) is using Docker Compose:

1. **Copy environment example and configure**
   ```bash
   cp .env.example .env
   # Edit .env with your Spotify credentials
   ```

2. **Start all services**
   ```bash
   docker compose up --build
   ```

   This will start:
   - **API**: Running on `http://127.0.0.1:3000`
   - **PostgreSQL**: Running on `127.0.0.1:5432`
   - **Redis**: Running on `127.0.0.1:6379`

3. **Stop services**
   ```bash
   docker compose down
   ```

4. **View logs**
   ```bash
   docker compose logs -f api
   docker compose logs -f postgres
   docker compose logs -f redis
   ```

### Docker Compose Services

#### API Service
- Image: `golang:1.22`
- Port: `3000`
- Auto-runs: `go run ./src/main.go`
- Depends on: PostgreSQL and Redis

#### PostgreSQL Service
- Image: `postgres:16`
- Port: `5432`
- Database: `clear_songs`
- Volume: `pgdata` (persistent)

#### Redis Service
- Image: `redis:7`
- Port: `6379`
- Persistence: Enabled (`appendonly yes`)

### Environment Variables in Compose

The compose file automatically configures:
- Database host points to `postgres` service
- Redis host points to `redis` service
- All services share a Docker network

Customize via `.env` file before running `docker compose up`.

### Required Spotify Permissions

The application requires the following Spotify scopes:

- `playlist-read-private` - Read private playlists
- `playlist-read-collaborative` - Read collaborative playlists
- `playlist-modify-public` - Modify public playlists
- `playlist-modify-private` - Modify private playlists
- `user-library-read` - Read user's saved tracks
- `user-library-modify` - Modify user's saved tracks
- `user-read-private` - Read user profile
- `user-read-email` - Read user email

## 📚 API Documentation

### Base URL

```
http://127.0.0.1:3000
```

### Interactive Documentation

Access the Swagger UI at: `http://127.0.0.1:3000/swagger/index.html`

---

## 🔐 Authentication Endpoints

### Login to Spotify

Initiates the OAuth flow to authenticate with Spotify.

**Endpoint:** `GET /auth/login`

**Response:**

- `302 Redirect` - Redirects to Spotify authentication page

**Example:**

```bash
curl -X GET "http://127.0.0.1:3000/auth/login"
```

### OAuth Callback

Handles the callback from Spotify after user authentication.

**Endpoint:** `GET /auth/callback`

**Query Parameters:**

- `code` (string, required) - Authorization code from Spotify

**Response:**

```json
{
  "status": "success",
  "message": "User authenticated"
}
```

### Check Authentication Status

Verifies if the user is currently authenticated.

**Endpoint:** `GET /auth/status`

**Response:**

```json
{
  "status": "success",
  "message": "User authenticated"
}
```

**Error Response:**

```json
{
  "status": "error",
  "message": "Unauthorized"
}
```

### Logout

Clears the current user session.

**Endpoint:** `POST /auth/logout`

**Response:**

```json
{
  "status": "success",
  "message": "User logged out"
}
```

---

## 🎵 Track Management Endpoints

### Get Track Summary

Returns a comprehensive summary of tracks organized by artist, with optional filtering.

**Endpoint:** `GET /track/summary`

**Query Parameters:**

- `min` (integer, optional) - Minimum track count filter
- `max` (integer, optional) - Maximum track count filter

**Response:**

```json
[
  {
    "id": "4NHQUGzhtTLFvgF5SZesLK",
    "name": "Radiohead",
    "count": 45
  },
  {
    "id": "1dfeR4HaWDbWqFHLkxsg1d",
    "name": "Queen",
    "count": 32
  }
]
```

**Example:**

```bash
# Get all artists
curl -X GET "http://127.0.0.1:3000/track/summary"

# Get artists with 10-50 tracks
curl -X GET "http://127.0.0.1:3000/track/summary?min=10&max=50"

# Get artists with more than 20 tracks
curl -X GET "http://127.0.0.1:3000/track/summary?min=20"
```

### Delete Tracks by Artist

Removes all tracks from a specific artist from your library.

**Endpoint:** `DELETE /track/artist/{id_artist}`

**Path Parameters:**

- `id_artist` (string, required) - Spotify Artist ID

**Response:**

```json
{
  "message": "Tracks deleted"
}
```

**Example:**

```bash
curl -X DELETE "http://127.0.0.1:3000/track/artist/4NHQUGzhtTLFvgF5SZesLK"
```

### Delete Tracks by Range

Removes tracks based on the number of songs per artist within a specified range.

**Endpoint:** `DELETE /track/range`

**Query Parameters:**

- `min` (integer, optional) - Minimum track count (artists with at least this many tracks)
- `max` (integer, optional) - Maximum track count (artists with at most this many tracks)

**Response:**

```json
{
  "message": "Tracks deleted"
}
```

**Examples:**

```bash
# Delete tracks from artists with exactly 1 track (likely singles)
curl -X DELETE "http://127.0.0.1:3000/track/range?min=1&max=1"

# Delete tracks from artists with more than 50 tracks
curl -X DELETE "http://127.0.0.1:3000/track/range?min=50"

# Delete tracks from artists with 5-15 tracks
curl -X DELETE "http://127.0.0.1:3000/track/range?min=5&max=15"
```

---

## 📋 Playlist Management Endpoints

### Delete All Playlist Tracks

Removes all tracks from a specified playlist while keeping the playlist structure intact.

**Endpoint:** `DELETE /playlist/tracks`

**Query Parameters:**

- `id` (string, required) - Spotify Playlist ID

**Response:**

```json
{
  "message": "Tracks deleted"
}
```

**Example:**

```bash
curl -X DELETE "http://127.0.0.1:3000/playlist/tracks?id=37i9dQZF1DXcBWIGoYBM5M"
```

### Delete Playlist Tracks and Remove from Library

Removes all tracks from both the specified playlist AND your personal library. Includes automatic backup to database.

**Endpoint:** `DELETE /playlist/tracks/all`

**Query Parameters:**

- `id` (string, required) - Spotify Playlist ID

**Response:**

```json
{
  "message": "Tracks deleted"
}
```

**Example:**

```bash
curl -X DELETE "http://127.0.0.1:3000/playlist/tracks/all?id=37i9dQZF1DXcBWIGoYBM5M"
```

**⚠️ Warning:** This operation removes tracks from your library permanently. Tracks are backed up to the database for recovery.

---

## 💿 Album Management Endpoints

### Convert Album to Individual Songs

Converts an album to individual songs in your library.

**Endpoint:** `POST /album/convert`

**Query Parameters:**

- `id_album` (string, required) - Spotify Album ID

**Response:**

```json
{
  "message": "Album converted to songs"
}
```

**Example:**

```bash
curl -X POST "http://127.0.0.1:3000/album/convert?id_album=4aawyAB9vmqN3uQ7FjRGTy"
```

---

## 🔧 Error Handling

The API uses standard HTTP status codes and returns detailed error messages:

### Common Error Responses

**400 Bad Request**

```json
{
  "message": "Playlist id is required"
}
```

**401 Unauthorized**

```json
{
  "status": "error",
  "message": "Unauthorized"
}
```

**500 Internal Server Error**

```json
{
  "message": "Error deleting tracks",
  "error": "detailed error description"
}
```

---

## 🎯 Usage Examples

### Complete Workflow Example

1. **Authenticate with Spotify**

   ```bash
   # Open browser and visit
   http://127.0.0.1:3000/auth/login
   ```

2. **Analyze your library**

   ```bash
   curl -X GET "http://127.0.0.1:3000/track/summary"
   ```

3. **Remove artists with only 1 track (cleanup singles)**

   ```bash
   curl -X DELETE "http://127.0.0.1:3000/track/range?min=1&max=1"
   ```

4. **Clean up a specific playlist**

   ```bash
   curl -X DELETE "http://127.0.0.1:3000/playlist/tracks?id=YOUR_PLAYLIST_ID"
   ```

5. **Verify changes**

   ```bash
   curl -X GET "http://127.0.0.1:3000/track/summary"
   ```

### Advanced Use Cases

**Remove duplicate artists (keep only artists with 10+ tracks):**

```bash
curl -X DELETE "http://127.0.0.1:3000/track/range?min=1&max=9"
```

**Clean library and specific playlist simultaneously:**

```bash
curl -X DELETE "http://127.0.0.1:3000/playlist/tracks/all?id=PLAYLIST_ID"
```

---

## 🏗️ Architecture & Performance

### Caching System

- **Smart Caching**: Automatic caching of Spotify API responses
- **Intelligent Invalidation**: Cache automatically updates after modifications
- **Performance Optimization**: Reduces API calls and improves response times

### Database Backup

- **Automatic Backup**: All deleted tracks are saved to PostgreSQL
- **Recovery Ready**: Easy restoration of accidentally deleted content
- **Data Integrity**: GORM ensures safe database operations

### Rate Limiting

The application respects Spotify API rate limits and implements proper pagination for large datasets.

---

## 🛠️ Development

### Project Structure

```
src/
├── cache/              # Cache management
├── constants/          # Application constants
├── controllers/        # HTTP handlers
├── database/          # Database configuration
├── docs/              # Swagger documentation
├── helpers/           # Utility helpers
├── middlewares/       # HTTP middlewares
├── models/            # Data models
├── routes/            # Route definitions
├── services/          # Business logic
└── utils/             # Utility functions
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o clear-songs src/main.go
```

### API Documentation Generation

```bash
swag init -g src/main.go
```

---

## ⚠️ Safety Considerations

### Backup System

- All track deletions are automatically backed up to PostgreSQL
- Recovery is possible through direct database access
- Consider implementing a recovery endpoint for easier restoration

### Rate Limiting

- Spotify API has rate limits - the application handles these gracefully
- Large operations are automatically paginated

### Recommendations

1. **Test First**: Always test with a small playlist or artist first
2. **Backup**: Although automatic backup is provided, consider exporting your library before major operations
3. **Verification**: Use the summary endpoint to verify changes
4. **Gradual Approach**: For large libraries, perform operations in smaller chunks

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [Spotify Web API](https://developer.spotify.com/documentation/web-api/) for providing comprehensive music data access
- [Gin Framework](https://gin-gonic.com/) for the excellent HTTP framework
- [GORM](https://gorm.io/) for elegant database operations

---

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/RubenPari/clear-songs/issues) page
2. Create a new issue with detailed information
3. Include logs and error messages when possible

**Happy music library management! 🎵**
