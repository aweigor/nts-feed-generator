package ntsclient

// Response is the top-level envelope returned by the API.
type EpisodesResponse struct {
	Metadata Metadata      `json:"metadata"`
	Results  []EpisodeInfo `json:"results"`
}

// Metadata wraps the pagination info for a response.
type Metadata struct {
	Resultset Resultset `json:"resultset"`
}

// Resultset describes pagination over the full result set.
type Resultset struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// Result is a single item in the results list.
type EpisodeInfo struct {
	ArticleType    string        `json:"article_type"`
	Title          string        `json:"title"`
	Artists        []Artist      `json:"artists"`
	Article        Article       `json:"article"`
	AudioSources   []AudioSource `json:"audio_sources"`
	Description    Description   `json:"description"`
	Image          Image         `json:"image"`
	RelatedEpisode *Episode      `json:"related_episode,omitempty"`
	LocalDate      string        `json:"local_date"`
	Location       string        `json:"location"`
	Genres         []Genre       `json:"genres"`
	Moods          []Mood        `json:"moods"`
	TrackUID       *string       `json:"track_uid"`
	Brand          *Brand        `json:"brand,omitempty"`
}

// Artist is a performer or contributor on a result.
type Artist struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Article points to the canonical page for a result.
type Article struct {
	Path string `json:"path"`
}

// AudioSource is a playable source for a result.
type AudioSource struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}

// Description is currently an empty object in the API but is modeled
// as a map so future keys do not break decoding.
type Description struct {
	// Fields are currently unknown; keep as a generic map.
	Fields map[string]any `json:"-"`
}

// Image holds the various resized image URLs for a result.
type Image struct {
	Large       string `json:"large"`
	MediumLarge string `json:"medium_large"`
	Medium      string `json:"medium"`
	Small       string `json:"small"`
	Thumb       string `json:"thumb"`
}

// Episode is a related episode reference.
type Episode struct {
	Path string `json:"path"`
}

// Genre is a musical genre tag.
type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Mood is a mood tag.
type Mood struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Brand describes the promotional brand block attached to some results.
type Brand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CTAURL      string `json:"cta_url"`
	CTALabel    string `json:"cta_label"`
}
