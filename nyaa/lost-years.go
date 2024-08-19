package nyaa

import (
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	lost_years_toorent = "lost_years-torrent"
	lost_years_anime   = "lost_years-anime"
	lost_years_path    = "LostYears"
)

func parseLostYears(i1 string) (*models.Query, error) {

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

	i1 = strings.ReplaceAll(i1, "|", "")
	i1 = strings.TrimSpace(i1)

	item := new(models.Query)
	item.TorrentFile = torrent
	globals = strings.Split(i1, "-")

	for i := len(globals) - 1; i >= 0; i-- {
		temp = strings.ReplaceAll(globals[i], " ", "")
		season, ep := analyze.ExtractSE(temp)
		if season != 0 && ep != 0 {
			item.Season = season
			item.Episodes = []float64{ep}
			continue
		}
		item.AnimeTitle = new(models.AnimeTitle)
		item.RomanjiTitle = analyze.CleanSymbols(strings.Join(globals[:i+1], "-"))
		item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
		break
	}

	return item, nil
}
