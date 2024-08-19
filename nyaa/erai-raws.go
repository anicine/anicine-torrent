package nyaa

import (
	"errors"
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	erai_raws_toorent = "erai_raws-torrent"
	erai_raws_anime   = "erai_raws-anime"
	erai_raws_path    = "erai-raws"
)

func parseEraiRaws(i1 string) (*models.Query, error) {
	var (
		torrent = new(models.TorrentFile)
		keys    []string
		globals []string
		temp    string
	)

	globals = analyze.ExtractBrackets(i1)

	for _, v1 := range globals {
		i1 = strings.Replace(i1, v1, "", 1)
		if temp = analyze.ExtractLanguage(v1); temp != "" {
			torrent.Sub = append(torrent.Sub, temp)
		} else {
			keys = append(keys, v1)
		}
	}

	globals = []string{}
	i1 = strings.TrimSpace(i1)

	for _, v1 := range keys {
		if temp = analyze.ExtractQuality(v1); temp != "" {
			torrent.Quality = temp
		} else {
			globals = append(globals, v1)
		}
	}

	globals = append(globals, analyze.ExtractParentheses(i1)...)
	keys = []string{}

	for _, v1 := range globals {
		i1 = strings.Replace(i1, v1, "", 1)
		if temp = analyze.ExtractACodec(v1); temp != "" {
			torrent.ACodec = temp
		} else {
			keys = append(keys, v1)
		}
	}

	globals = []string{}
	for _, v1 := range keys {
		if temp = analyze.ExtractVCodec(v1); temp != "" {
			torrent.VCodec = temp
		} else {
			globals = append(globals, v1)
		}
	}

	for _, v1 := range globals {
		v1 = strings.ToLower(strings.TrimSpace(v1))
		if strings.Contains(v1, "raw") {
			continue
		}
		if len(v1) == 10 {
			torrent.CRC = analyze.CleanSymbols(v1)
			break
		}
	}

	item := new(models.Query)
	item.TorrentFile = torrent
	globals = strings.Split(i1, "-")

	if len(globals) < 2 {
		return nil, errors.New("cannot parse the episodes")
	}

	for i := len(globals) - 1; i >= 0; i-- {
		temp = strings.ReplaceAll(globals[i], " ", "")
		temp = strings.ReplaceAll(temp, "_", "~")
		if ok := analyze.IsValidEpStr(temp); ok && item.Episodes == nil {
			eps := analyze.ExtractFloatsWithRanges(temp)
			item.Batch = true
			item.Episodes = eps
			continue
		}
		temp = strings.Split(temp, "v")[0]
		if item.Episodes == nil {
			if ok := analyze.IsValidEpStr(temp); ok {
				eps := analyze.ExtractFloatsWithRanges(temp)
				if len(eps) > 1 {
					item.Batch = true
				}
				item.Episodes = eps
				continue
			}
			if strings.Contains(strings.ToLower(temp), "movie") {
				item.Movie = true
				item.Episodes = []float64{1}
				continue
			}

		}
		if i < len(globals)-1 {
			item.AnimeTitle = new(models.AnimeTitle)
			item.RomanjiTitle = strings.Join(globals[:i+1], "-")
			item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
			break
		} else {
			return nil, errors.New("cannot parse the anime name")
		}
	}

	return item, nil
}
