package nyaa

import (
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	beatrice_raws_toorent = "beatrice_raws-torrent"
	beatrice_raws_anime   = "beatrice_raws-anime"
	beatrice_raws_path    = "DJATOM"
)

func parseBeatriceRaws(i1 string) (*models.Query, error) {
	var (
		torrent       = new(models.TorrentFile)
		keys, globals []string
		temp          string
	)

	keys = analyze.ExtractParentheses(i1)
	keys = append(keys, analyze.ExtractBrackets(i1)...)

	for _, v1 := range keys {
		v1 = analyze.CleanSymbols(v1)
		globals = append(globals, strings.Split(v1, " ")...)
	}

	keys = []string{}
	for _, v1 := range globals {
		if temp = analyze.ExtractVCodec(v1); temp != "" {
			torrent.VCodec = temp
		} else {
			keys = append(keys, v1)
		}
	}

	for _, v1 := range keys {
		i1 = strings.Replace(i1, v1, "", 1)
		if temp = analyze.ExtractACodec(v1); temp != "" {
			torrent.ACodec = temp
		} else {
			globals = append(globals, v1)
		}

	}

	keys = []string{}
	for _, v1 := range globals {
		i1 = strings.Replace(i1, v1, "", 1)
		if temp = analyze.ExtractQuality(v1); temp != "" {
			torrent.Quality = temp
			continue
		} else {
			keys = append(keys, v1)
		}
	}

	i1 = analyze.CleanSymbols(strings.TrimSpace(i1))

	if strings.Contains(i1, `/`) {
		globals = strings.Split(i1, `/`)
	} else if strings.Contains(i1, `\`) {
		globals = strings.Split(i1, `\`)
	} else if strings.Contains(i1, `|`) {
		globals = strings.Split(i1, `|`)
	} else {
		globals = []string{}
	}

	item := new(models.Query)
	item.TorrentFile = torrent
	item.AnimeTitle = new(models.AnimeTitle)

	if len(globals) > 0 {
		if len(globals) == 2 {
			item.AnimeTitle.RomanjiTitle = strings.TrimSpace(globals[0])
			item.AnimeTitle.EnglishTitle = strings.TrimSpace(globals[1])
		} else {
			item.AnimeTitle.RomanjiTitle = strings.TrimSpace(globals[0])
		}
	} else {
		item.AnimeTitle.RomanjiTitle = i1
	}

	return item, nil
}
