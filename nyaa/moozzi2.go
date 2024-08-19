package nyaa

import (
	"errors"
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	moozzi2_toorent = "moozzi2-torrent"
	moozzi2_anime   = "moozzi2-anime"
	moozzi2_path    = "Moozzi2"
)

func parseMoozzi2(i1 string) (*models.Query, error) {

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

	keys = []string{}
	for _, v1 := range globals {
		if temp = analyze.ExtractQuality(v1); temp != "" {
			torrent.Quality = temp
			continue
		} else {
			keys = append(keys, v1)
		}
	}

	globals = []string{}
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
			break
		}
	}

	item := new(models.Query)
	item.TorrentFile = torrent

	i1 = strings.TrimSpace(i1)
	globals = strings.Split(i1, "-")

	if len(globals) == 1 {
		globals = strings.Split(i1, "BD")
	}

	if len(globals) == 1 {
		return nil, errors.New("cannot parse the anime name and type")
	}

	temp = strings.TrimSpace(globals[0])
	temp = strings.TrimSuffix(temp, "BD")
	temp = strings.TrimSpace(temp)

	item.AnimeTitle = new(models.AnimeTitle)
	item.AnimeTitle.RomanjiTitle = globals[0][:len(strings.Split(temp, ""))]

	keys = append([]string{}, globals[1:]...)

	for _, key := range keys {
		temp = strings.ToLower(key)
		if strings.Contains(temp, "movie") {
			item.Movie = true
			break
		}
		if strings.Contains(temp, "tv") {
			item.Batch = true
			break
		}
	}

	return item, nil
}
