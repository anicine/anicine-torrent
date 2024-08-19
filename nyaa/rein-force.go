package nyaa

import (
	"fmt"
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	rein_force_toorent = "rein_force-torrent"
	rein_force_anime   = "rein_force-anime"
	rein_force_path    = "ReinForce"
)

func parseReinForce(i1 string) (*models.Query, error) {

	var (
		torrent = new(models.TorrentFile)
		keys    []string
		globals []string
		temp    string
	)
	keys = analyze.ExtractParentheses(i1)
	keys = append(keys, analyze.ExtractBrackets(i1)...)

	for _, v1 := range keys {
		i1 = strings.Replace(i1, v1, "", 1)
		v1 = analyze.CleanSymbols(v1)
		globals = append(globals, strings.Split(v1, " ")...)
	}

	fmt.Println(globals)

	keys = []string{}
	for _, v1 := range globals {
		if temp = analyze.ExtractQuality(v1); temp != "" {
			torrent.Quality = temp
			continue
		} else {
			keys = append(keys, v1)
		}
	}

	for _, v1 := range keys {
		if temp = analyze.ExtractVCodec(v1); temp != "" {
			torrent.VCodec = temp
			continue
		} else {
			globals = append(globals, v1)
		}
	}

	for _, v1 := range globals {
		if temp = analyze.ExtractACodec(v1); temp != "" {
			torrent.ACodec = temp
			continue
		} else {
			keys = append(keys, v1)
		}
	}

	i1 = strings.TrimSpace(i1)

	item := new(models.Query)
	item.TorrentFile = torrent
	item.AnimeTitle = new(models.AnimeTitle)
	item.AnimeTitle.RomanjiTitle = i1

	return item, nil
}
