package nyaa

import (
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	yameii_toorent = "yameii-torrent"
	yameii_anime   = "yameii-anime"
	yameii_path    = "Yameii"
)

func parseYameii(i1 string) (*models.Query, error) {

	var (
		torrent = new(models.TorrentFile)
		keys    []string
		globals []string
		temp    string
	)

	globals = analyze.ExtractParentheses(i1)
	globals = append(globals, analyze.ExtractBrackets(i1)...)

	for _, v1 := range globals {
		i1 = strings.Replace(i1, v1, "", 1)
		if temp = analyze.ExtractQuality(v1); temp != "" {
			torrent.Quality = temp
			continue
		} else {
			keys = append(keys, v1)
		}
		if len(torrent.Dub) == 0 {
			if strings.Contains(i1, "English Dub") || strings.Contains(i1, "ENG Dub") {
				torrent.Dub = append(torrent.Dub, "ENG")
				torrent.Sub = append(torrent.Sub, "ENG")
			}
		}
	}

	i1 = strings.TrimSpace(i1)
	for _, v1 := range keys {
		v1 = strings.ToLower(strings.TrimSpace(v1))
		if strings.Contains(v1, "yameii") {
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

	return item, nil
}
