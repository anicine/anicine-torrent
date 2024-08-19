package models

type Query struct {
	*AnimeTitle
	*TorrentFile
	Season       int
	Movie, Batch bool
	Episodes     []float64
}
