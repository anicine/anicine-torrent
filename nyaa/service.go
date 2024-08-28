package nyaa

import (
	"log/slog"

	"github.com/anacrolix/torrent"
	"go.mongodb.org/mongo-driver/mongo"
)

type Nyaa struct {
	log      *slog.Logger
	db       *mongo.Database
	client   *torrent.Client
	record   *record
	trackers []string
}

func NewNyaa(i1 *slog.Logger, i2 *mongo.Database, i3 *torrent.Client, i4 []string) *Nyaa {
	return &Nyaa{
		log:      i1,
		db:       i2,
		client:   i3,
		record:   new(record),
		trackers: i4,
	}
}
