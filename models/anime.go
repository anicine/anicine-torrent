package models

type AnimeTitle struct {
	OriginalTitle string
	RomanjiTitle  string
	EnglishTitle  string
}

type Anime struct {
	*AnimeTitle
	Type     string
	Season   int
	Batch    AnimeBatch
	Episodes AnimeEpisode
}

type AnimeEpisode struct {
	*Torrent
	Episode float64
}

type AnimeBatch struct {
	*TorrentFile
	Episodes  []map[float64][]*TorrentData
	OtherData []*TorrentData
}
