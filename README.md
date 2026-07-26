1. System Overview
A single Go application with two independent modules that can be separated later.

2. Module Architecture
text
┌─────────────────────────────────────────────────────────┐
│                    Application                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────┐         ┌─────────────────────┐  │
│  │                 │         │                     │  │
│  │   Feed Module   │────────▶│ Tracklist Module    │  │
│  │                 │  calls  │ (NTS only)          │  │
│  │  - RSS/JSON     │         │                     │  │
│  │  - HTTP server  │         │  - Fetches NTS      │  │
│  │  - Scheduler    │         │    tracklists       │  │
│  │                 │         │  - Simple cache     │  │
│  └─────────────────┘         └─────────────────────┘  │
│                                                         │
└─────────────────────────────────────────────────────────┘
3. Module Responsibilities
Feed Module
Fetches shows/episodes from NTS API

Generates RSS and JSON feeds

Calls Tracklist Module to get tracklist URLs

Serves feeds via HTTP

Tracklist Module (NTS Only)
Fetches tracklists from NTS API

Simple cache for tracklists

Returns tracklist data or URL

4. Simple Flow
text
1. HTTP request for feed
2. Feed module builds episode list
3. For each episode:
   - Calls tracklist module: GetTracklist(episodeID)
   - Tracklist module checks NTS API
   - Returns tracklist data (or empty)
4. Feed module includes tracklist URL in RSS/JSON
5. Returns feed to client
5. Module Interface
go
// Tracklist module exposes this interface
type TracklistProvider interface {
    GetTracklist(episodeID string) (*Tracklist, error)
}
6. Module Structure
text
internal/
├── feed/
│   ├── generator.go    # RSS/JSON generation
│   ├── handlers.go     # HTTP endpoints
│   └── scheduler.go    # Periodic refresh
│
├── tracklist/           # ← NTS only, separate module
│   ├── provider.go     # Implements interface
│   ├── client.go       # Calls NTS API
│   └── cache.go        # Simple in-memory cache
│
└── nts/
    └── client.go       # Shared NTS API client

8. What Tracklist Module Does Have
✅ NTS API client for tracklist endpoint

✅ Simple in-memory cache (TTL: 1 hour)

✅ Error handling for missing tracklists

✅ Clean interface for feed module

9. Data Format
Tracklist Module Returns:
go
type Tracklist struct {
    EpisodeID   string
    Source      string  // Always "nts"
    Tracks      []Track
    FetchedAt   time.Time
}
Feed Module Uses:
go
// RSS includes tracklist_url
// JSON includes tracklist_url
// Both point to: /api/episodes/{id}/tracklist
10. Future Separation
Tracklist Module Can Be Extracted Because:
Has clear interface

No dependencies on feed module

Only depends on NTS API

Self-contained caching

To Make It a Separate Service Later:
text
Current: internal/tracklist/ (local package)
Future: tracklist-service/ (separate repo)
11. Configuration
yaml
nts:
  api_url: "https://www.nts.live/api/v2/"

tracklist:
  cache_ttl: "1h"
  enabled: true

feed:
  output_dir: "./feeds"
  include_tracklist_url: true


Feed Structure
1. Feed Categories (Based on NTS API)
1.1 Per-Show Feeds
text
/feeds/shows/{show_id}.{xml|json}
One feed per show

Contains all episodes for that show

Episodes sorted newest first

1.2 Latest Episodes Feed
text
/feeds/latest.{xml|json}
Aggregates latest episodes across all shows

Configurable limit (e.g., 20 episodes)

Sorted by air date descending

1.3 Live Channel Feeds
text
/feeds/live/1.{xml|json}
/feeds/live/2.{xml|json}
Current show on NTS 1 or NTS 2

Updates when show changes

Single item per feed

1.4 Infinite Mixtapes Feeds
text
/feeds/mixtape/{mixtape_id}.{xml|json}
Each mixtape gets its own feed

Continuous stream, updated regularly

Examples: "slow-focus", "poolside", "4-to-the-floor"

1.5 Genre Feeds
text
/feeds/genre/{genre}.{xml|json}
Filter episodes by genre

Examples: "jazz", "techno", "afrobeat", "disco"

Episodes tagged with that genre

1.6 Location Feeds
text
/feeds/location/{city}.{xml|json}
Shows from specific city

Examples: "london", "new-york", "tokyo"

Episodes from shows broadcasting from that city

1.7 Artist Spotlight Feed
text
/feeds/artist/{artist_name}.{xml|json}
Search episodes where artist appears

Dynamically generated on request

Uses NTS tracklist search

1.8 Staff Picks/Curated Feeds
text
/feeds/curated/staff-picks.{xml|json}
/feeds/curated/recommended.{xml|json}
Manually curated episodes

Updated periodically

Admin-controlled selection

2. URL Structure
Base Pattern
text
GET /feeds/{type}/{identifier}.{format}

Where:
- {type}: shows, latest, live, mixtape, genre, location, artist, curated
- {identifier}: show_id, city, genre, artist_name, etc.
- {format}: xml (RSS) or json
Examples
text
# Show feed
GET /feeds/shows/london-jazz.xml
GET /feeds/shows/london-jazz.json

# Latest episodes
GET /feeds/latest.xml
GET /feeds/latest.json

# Live channel
GET /feeds/live/1.xml
GET /feeds/live/2.xml

# Mixtape
GET /feeds/mixtape/poolside.xml
GET /feeds/mixtape/poolside.json

# Genre
GET /feeds/genre/jazz.xml
GET /feeds/genre/techno.json

# Location
GET /feeds/location/tokyo.xml
GET /feeds/location/new-york.json

# Artist
GET /feeds/artist/aphex-twin.xml
GET /feeds/artist/bjork.json

# Curated
GET /feeds/curated/staff-picks.xml
GET /feeds/curated/recommended.json
3. Feed Index
List All Available Feeds
text
GET /feeds
Response:

json
{
  "feeds": [
    {
      "type": "show",
      "name": "London Jazz",
      "url": "/feeds/shows/london-jazz.xml",
      "json_url": "/feeds/shows/london-jazz.json"
    },
    {
      "type": "show",
      "name": "Global Roots",
      "url": "/feeds/shows/global-roots.xml"
    },
    {
      "type": "latest",
      "name": "Latest Episodes",
      "url": "/feeds/latest.xml"
    },
    {
      "type": "live",
      "name": "NTS 1 Live",
      "url": "/feeds/live/1.xml"
    },
    {
      "type": "mixtape",
      "name": "Poolside",
      "url": "/feeds/mixtape/poolside.xml"
    }
  ]
}
4. Feed Item Structure
Common Fields for All Feed Types
xml
<item>
    <title>Episode Title</title>
    <link>https://www.nts.live/episodes/{id}</link>
    <description>Episode description</description>
    <pubDate>2024-01-15T20:00:00Z</pubDate>
    <enclosure url="{audio_url}" type="audio/mpeg" length="0"/>
    <guid>episode-123</guid>
    
    <!-- Custom fields -->
    <nts:show>Show Name</nts:show>
    <nts:air_date>2024-01-15</nts:air_date>
    <nts:duration>120</nts:duration>
    <nts:tracklist_url>/api/episodes/123/tracks</nts:tracklist_url>
</item>
5. Configuration
yaml
feed:
  output_dir: "./feeds"
  formats: ["xml", "json"]  # RSS and JSON
  
  latest:
    limit: 20
    cache_ttl: "15m"
  
  live:
    update_interval: "5m"
  
  mixtape:
    update_interval: "15m"
  
  genre:
    enabled: true
  
  location:
    enabled: true
  
  artist:
    cache_ttl: "1h"
    limit: 50
  
  curated:
    enabled: true
    manual_list: ["staff-picks", "recommended"]
6. Module Flow
text
┌──────────────────────────────────────────────────────┐
│                   Feed Module                        │
├──────────────────────────────────────────────────────┤
│                                                      │
│  HTTP Request ──► Router                            │
│                      │                              │
│          ┌───────────┴──────────────┐              │
│          │                          │              │
│          ▼                          ▼              │
│    Show Handler              Latest Handler         │
│    Live Handler             Mixtape Handler         │
│    Genre Handler            Location Handler        │
│    Artist Handler           Curated Handler         │
│          │                          │              │
│          └───────────┬──────────────┘              │
│                      │                              │
│                      ▼                              │
│              Fetch Episodes                        │
│                      │                              │
│                      ▼                              │
│         Tracklist Module (NTS)                    │
│                      │                              │
│                      ▼                              │
│         Generate RSS/JSON                         │
│                      │                              │
│                      ▼                              │
│              HTTP Response                        │
│                                                      │
└──────────────────────────────────────────────────────┘
7. Summary
Feed Types:

Show feeds: One per show

Latest: All shows aggregated

Live: Current broadcasts (NTS 1, NTS 2)

Mixtapes: 24/7 thematic streams

Genre: Filtered by music genre

Location: Filtered by city

Artist: Search by artist name

Curated: Manual selections

Format: RSS 2.0 (.xml) and JSON Feed 1.1 (.json)
