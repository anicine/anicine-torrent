package nyaa

import (
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	asw_toorent = "asw-torrent"
	asw_anime   = "asw-anime"
	asw_path    = "AkihitoSubsWeeklies"
)

func parseASW(i1 string) (*models.Query, error) {
	var (
		torrent       = new(models.TorrentFile)
		keys, globals []string
		temp          string
		batch         bool
	)

	keys = analyze.ExtractParentheses(i1)
	keys = append(keys, analyze.ExtractBrackets(i1)...)

	for _, v1 := range keys {
		v1 = analyze.CleanSymbols(v1)
		globals = append(globals, strings.Split(v1, " ")...)
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

	i1 = strings.TrimSpace(i1)

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
		if temp = analyze.ExtractVCodec(v1); temp != "" {
			torrent.VCodec = temp
		} else {
			keys = append(keys, v1)
		}
	}

	for _, v1 := range keys {
		v1 = strings.ToLower(strings.TrimSpace(v1))
		if strings.Contains(v1, "batch") {
			batch = true
			continue
		}
		if strings.Contains(v1, "asw") {
			continue
		}
		if len(v1) == 8 {
			torrent.CRC = analyze.CleanSymbols(v1)
		}
	}

	item := new(models.Query)
	item.TorrentFile = torrent

	if !batch {
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
			temp = strings.Split(temp, "v")[0]
			if ok := analyze.IsValidEpStr(temp); ok && item.Episodes == nil {
				eps := analyze.ExtractFloatsWithRanges(temp)
				if len(eps) > 1 {
					item.Batch = true
				}
				item.Episodes = eps
				continue
			}
			if item.Episodes == nil {
				item.Movie = true
				item.Episodes = []float64{1}
			}

			item.AnimeTitle = new(models.AnimeTitle)
			item.RomanjiTitle = analyze.CleanSymbols(strings.Join(globals[:i+1], "-"))
			item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
			break
		}
	} else {
		item.AnimeTitle = new(models.AnimeTitle)
		item.RomanjiTitle = strings.TrimSpace(analyze.CleanSymbols(i1))
	}

	return item, nil
}
