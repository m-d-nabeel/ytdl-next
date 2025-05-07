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
	FormatID     string      `json:"format_id"`
	Resolution   string      `json:"resolution"`
	Ext          string      `json:"ext"`
	Filesize     int64       `json:"filesize"`
	Note         string      `json:"format_note"`
	Acodec       string      `json:"acodec"`
	Vcodec       string      `json:"vcodec"`
	IsCompatible bool        `json:"is_compatible,omitempty"` // Flag for cross-platform compatibility
	Quality      interface{} `json:"quality,omitempty"`       // Can be string or number in yt-dlp output
	QualityLabel string      `json:"-"`                       // Our computed quality label (not from JSON)
}

// Add helper methods to easily check audio/video availability
func (f *YTFormat) HasAudio() bool {
	return f.Acodec != "none"
}

func (f *YTFormat) HasVideo() bool {
	return f.Vcodec != "none"
}

// IsUniversallyCompatible returns true if the format uses widely supported codecs
func (f *YTFormat) IsUniversallyCompatible() bool {
	// Common universally compatible video codecs (h.264/AVC)
	universalVideoCodecs := map[string]bool{
		"avc1":        true,
		"h264":        true,
		"avc1.42e01e": true, // Baseline profile
		"avc1.4d401e": true, // Main profile
		"avc1.64001e": true, // High profile
		"avc1.640028": true, // High profile
	}

	// Common universally compatible audio codecs (AAC, MP3)
	universalAudioCodecs := map[string]bool{
		"aac":       true,
		"mp4a":      true,
		"mp4a.40.2": true, // AAC-LC
		"mp3":       true,
		"mp4a.40.5": true, // HE-AAC
	}

	// For video formats, check video codec compatibility
	if f.HasVideo() {
		// Check for common h.264 variant prefixes
		for prefix := range universalVideoCodecs {
			if len(f.Vcodec) >= len(prefix) && f.Vcodec[:len(prefix)] == prefix {
				// If it also has audio, check audio compatibility too
				if f.HasAudio() {
					for aPrefix := range universalAudioCodecs {
						if len(f.Acodec) >= len(aPrefix) && f.Acodec[:len(aPrefix)] == aPrefix {
							return true
						}
					}
					// Video compatible but audio isn't
					return false
				}
				// Video-only and compatible
				return true
			}
		}
		// Video codec not compatible
		return false
	}

	// For audio-only formats
	if f.HasAudio() {
		for prefix := range universalAudioCodecs {
			if len(f.Acodec) >= len(prefix) && f.Acodec[:len(prefix)] == prefix {
				return true
			}
		}
	}

	return false
}

// GetFormatType returns a descriptive type of the format
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

// IsVP9 returns true if the video codec is VP9 (widely supported but not universal)
func (f *YTFormat) IsVP9() bool {
	return f.HasVideo() && (f.Vcodec == "vp9" ||
		(len(f.Vcodec) >= 3 && f.Vcodec[:3] == "vp9"))
}

// IsAV1 returns true if the video codec is AV1 (newer, less compatible)
func (f *YTFormat) IsAV1() bool {
	return f.HasVideo() && (f.Vcodec == "av01" ||
		(len(f.Vcodec) >= 4 && f.Vcodec[:4] == "av01"))
}
