package main

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
}
