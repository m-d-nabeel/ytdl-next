package api

import (
	"os/exec"
)

// VideoQuality represents available video quality options
type VideoQuality string

const (
	QualityBest  VideoQuality = "best"
	Quality2160p VideoQuality = "2160"
	Quality1440p VideoQuality = "1440"
	Quality1080p VideoQuality = "1080"
	Quality720p  VideoQuality = "720"
	Quality480p  VideoQuality = "480"
	Quality360p  VideoQuality = "360"
)

// AudioQuality represents available audio quality options in kbps
type AudioQuality string

const (
	AudioQualityBest AudioQuality = "best"
	AudioQuality320  AudioQuality = "320"
	AudioQuality256  AudioQuality = "256"
	AudioQuality192  AudioQuality = "192"
	AudioQuality128  AudioQuality = "128"
	AudioQuality96   AudioQuality = "96"
	AudioQuality64   AudioQuality = "64"
)

// AudioFormat represents available audio formats
type AudioFormat string

const (
	AudioFormatMP3  AudioFormat = "mp3"
	AudioFormatM4A  AudioFormat = "m4a"
	AudioFormatWAV  AudioFormat = "wav"
	AudioFormatOPUS AudioFormat = "opus"
	AudioFormatAAC  AudioFormat = "aac"
)

// VideoFormat represents available video formats
type VideoFormat string

const (
	VideoFormatMP4  VideoFormat = "mp4"
	VideoFormatWEBM VideoFormat = "webm"
	VideoFormatMKV  VideoFormat = "mkv"
)

// MediaOptions represents common options for media download
type MediaOptions struct {
	IncludeEnglishSubtitles bool // Download English subtitles if available
	IncludeAutoSubtitles    bool // Download auto-generated English subtitles if available
}

// CmdInterface defines the interface for command execution
type CmdInterface interface {
	GetVideoInfoCmd(url string) *exec.Cmd
	GetSimpleDownloadCmd(url string) *exec.Cmd
	GetBestVideoCmd(url string, options *MediaOptions) *exec.Cmd
	GetVideoWithQualityCmd(url string, quality VideoQuality, format VideoFormat, options *MediaOptions) *exec.Cmd
	GetAudioOnlyCmd(url string, quality AudioQuality, format AudioFormat) *exec.Cmd
	GetMediaByFormatID(url, formatID string) *exec.Cmd
}

// GetVideoInfoCmd returns a command to fetch video information with thumbnail
func GetVideoInfoCmd(url string) *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args,
		"-j",                // Output as JSON
		"--write-thumbnail", // Get thumbnail URL/info
		"--skip-download",   // Don't download the video
		"--no-playlist",     // Don't process playlists
		url,
	)
	return cmd
}

func GetMediaByFormatID(url, formatID string) *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args,
		"-f", formatID, // Specific format ID
		"--no-playlist",     // Don't process playlists
		"--write-thumbnail", // Download thumbnail
		url,
	)
	return cmd
}

func GetSimpleDownloadCmd(url string) *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args,
		"-f", "best", // Best quality single file (no merging)
		"--no-playlist",
		"--write-thumbnail",
		url,
	)
	return cmd
}

// addSubtitleOptions adds subtitle options to the command if enabled
func addSubtitleOptions(cmd *exec.Cmd, options *MediaOptions) {
	if options != nil {
		if options.IncludeEnglishSubtitles {
			cmd.Args = append(cmd.Args,
				"--write-sub",
				"--sub-lang", "en",
			)
		}
		if options.IncludeAutoSubtitles {
			cmd.Args = append(cmd.Args,
				"--write-auto-sub",
				"--sub-lang", "en",
			)
		}
	}
}

// GetBestVideoCmd returns a command to download the best quality video
func GetBestVideoCmd(url string, options *MediaOptions) *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args,
		"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best", // Best quality MP4
		"--merge-output-format", "mp4", // Ensure MP4 output
		"--write-thumbnail", // Download thumbnail
		"--no-playlist",     // Don't process playlists
	)

	addSubtitleOptions(cmd, options)
	cmd.Args = append(cmd.Args, url)
	return cmd
}

// GetVideoWithQualityCmd returns a command to download video with specific quality and format
func GetVideoWithQualityCmd(url string, quality VideoQuality, format VideoFormat, options *MediaOptions) *exec.Cmd {
	cmd := exec.Command("yt-dlp")

	formatString := ""
	if quality == QualityBest {
		formatString = "bestvideo[ext=" + string(format) + "]+bestaudio/best[ext=" + string(format) + "]/best"
	} else {
		formatString = "bestvideo[height<=?" + string(quality) + "][ext=" + string(format) + "]+bestaudio/best[height<=?" + string(quality) + "][ext=" + string(format) + "]/best"
	}

	cmd.Args = append(cmd.Args,
		"-f", formatString,
		"--merge-output-format", string(format),
		"--write-thumbnail",
		"--no-playlist",
	)

	addSubtitleOptions(cmd, options)
	cmd.Args = append(cmd.Args, url)
	return cmd
}

// GetAudioOnlyCmd returns a command to download audio only with specific quality and format
func GetAudioOnlyCmd(url string, quality AudioQuality, format AudioFormat) *exec.Cmd {
	cmd := exec.Command("yt-dlp")
	cmd.Args = append(cmd.Args,
		"-f", "bestaudio/best", // Select best audio
		"-x", // Extract audio
		"--audio-format", string(format),
		"--audio-quality", string(quality),
		"--no-playlist",
		url,
	)
	return cmd
}
