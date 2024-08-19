package nyaa

import (
	"fmt"
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	ohys_raws_toorent = "ohys_raws-torrent"
	ohys_raws_anime   = "ohys_raws-anime"
	ohys_raws_path    = "ohys"
)

func parseOhysRaws(i1 string) (*models.Query, error) {

	var (
		torrent       = new(models.TorrentFile)
		keys, globals []string
		temp          string
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

	i1 = strings.Split(i1, "|")[0]
	i1 = strings.TrimSpace(i1)

	item := new(models.Query)
	item.TorrentFile = torrent
	globals = strings.Split(i1, "-")

	for i := len(globals) - 1; i >= 0; i-- {
		temp = strings.ReplaceAll(globals[i], " ", "")
		temp = strings.Split(temp, "v")[0]
		fmt.Println(temp)
		if ok := analyze.IsValidEpStr(temp); ok && item.Episodes == nil {
			eps := analyze.ExtractFloatsWithRanges(temp)
			if len(eps) > 1 {
				item.Batch = true
			}
			item.Episodes = eps
			continue
		}

		if item.Episodes == nil {
			if analyze.HasExtension(i1) {
				item.Movie = true
				item.Episodes = []float64{1}
			} else {
				item.Batch = true
			}
		}

		item.AnimeTitle = new(models.AnimeTitle)
		item.RomanjiTitle = analyze.CleanSymbols(strings.Join(globals[:i+1], "-"))
		item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
		break
	}

	return item, nil
}
