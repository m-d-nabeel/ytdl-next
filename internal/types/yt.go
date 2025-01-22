package types

type YTMediaInfo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Duration  int        `json:"duration"`
	Uploader  string     `json:"uploader"`
	ViewCount int        `json:"view_count"`
	LikeCount int        `json:"like_count"`
	Formats   []YTFormat `json:"formats"`
}

type YTFormat struct {
	FormatID   string `json:"format_id"`
	Resolution string `json:"resolution"`
	Ext        string `json:"ext"`
	Filesize   int64  `json:"filesize"`
	Note       string `json:"format_note"`
	Acodec     string `json:"acodec"` // Add this field
	Vcodec     string `json:"vcodec"` // Add this field
}

// Add helper methods to easily check audio/video availability
func (f *YTFormat) HasAudio() bool {
	return f.Acodec != "none"
}

func (f *YTFormat) HasVideo() bool {
	return f.Vcodec != "none"
}

// Optional: Add a method to get format type
func (f *YTFormat) GetFormatType() string {
	switch {
	case f.HasVideo() && f.HasAudio():
		return "video+audio"
	case f.HasVideo():
		return "video-only"
	case f.HasAudio():
		return "audio-only"
	default:
		return "unknown"
	}
}
