package models

type Torrent struct {
	ID string `json:"_id" bson:"_id"`
	*TorrentData
	Files []TorrentFile
}

type TorrentData struct {
	NyaaID int64   `json:"NyaaID" bson:"NyaaID"`
	Name   string  `json:"Name" bson:"Name"`
	Size   float64 `json:"Size" bson:"Size"`
	Date   int64   `json:"Date" bson:"Date"`
	*TorrentInfo
}

func (s *TorrentData) Clear() {
	s.NyaaID = 0
	s.Name = ""
	s.Size = 0
	s.Date = 0
	if s.TorrentInfo != nil {
		s.Seeders = 0
		s.Leechers = 0
		s.Hash = ""
	}
}

type TorrentInfo struct {
	Seeders  int    `json:"Seeders" bson:"Seeders"`
	Leechers int    `json:"Leechers" bson:"Leechers"`
	Hash     string `json:"Hash" bson:"Hash"`
	Magnet   string `json:"Magnet" bson:"Magnet"`
}

type TorrentFile struct {
	Size    int64    `json:"Size" bson:"Size"`
	Path    string   `json:"Path" bson:"Path"`
	CRC     string   `json:"CRC" bson:"CRC"`
	Quality string   `json:"Quality" bson:"Quality"`
	ACodec  string   `json:"ACodec" bson:"ACodec"`
	VCodec  string   `json:"VCodec" bson:"VCodec"`
	Sub     []string `json:"Sub" bson:"Sub"`
	Dub     []string `json:"Dub" bson:"Dub"`
}
