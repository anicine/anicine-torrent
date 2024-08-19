package nyaa

import (
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	ember_toorent = "ember-torrent"
	ember_anime   = "ember-anime"
	ember_path    = "ember_encodes"
)

func parseEmber(i1 string) (*models.Query, error) {

	var (
		torrent       = new(models.TorrentFile)
		keys, globals []string
		movie, batch  bool
		temp          string
	)

	globals = analyze.ExtractParentheses(i1)
	globals = append(globals, analyze.ExtractBrackets(i1)...)

	for _, v1 := range globals {
		i1 = strings.Replace(i1, v1, "", 1)
		if temp = analyze.ExtractQuality(v1); temp != "" {
			torrent.Quality = temp
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

	keys = []string{}
	for _, v1 := range globals {
		if temp = analyze.ExtractACodec(v1); temp != "" {
			torrent.ACodec = temp
		} else {
			keys = append(keys, v1)
		}
	}

	for _, v1 := range keys {
		v1 = strings.ToLower(strings.TrimSpace(v1))
		if strings.Contains(v1, "batch") {
			batch = true
			movie = false
		}
		if strings.Contains(v1, "movie") {
			movie = true
		}
	}

	i1 = strings.TrimSpace(i1)

	item := new(models.Query)
	item.TorrentFile = torrent
	if !batch {
		if !movie {
			globals = strings.Split(i1, "-")
			for i := len(globals) - 1; i >= 0; i-- {
				temp = strings.ReplaceAll(globals[i], " ", "")
				temp = strings.ReplaceAll(temp, "_", "~")
				season, ep := analyze.ExtractSE(temp)
				if season != 0 && ep != 0 {
					item.Season = season
					item.Episodes = []float64{ep}
					continue
				}
				temp = strings.Split(temp, "|")[0]
				temp = strings.Split(temp, "v")[0]
				if ok := analyze.IsValidEpStr(temp); ok && item.Episodes == nil {
					eps := analyze.ExtractFloatsWithRanges(temp)
					item.Episodes = eps
					continue
				}

				item.AnimeTitle = new(models.AnimeTitle)
				item.RomanjiTitle = analyze.CleanSymbols(strings.Join(globals[:i+1], "-"))
				item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
				break
			}
		} else {
			item.Movie = true
			item.Episodes = []float64{1}
			item.AnimeTitle = new(models.AnimeTitle)
			item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(analyze.CleanSymbols(i1)))
		}
	} else {
		item.Batch = true
		item.AnimeTitle = new(models.AnimeTitle)
		item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(analyze.CleanSymbols(i1)))
	}

	return item, nil
}
