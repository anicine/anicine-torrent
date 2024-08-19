package nyaa

import (
	"errors"
	"strings"

	"github.com/anicine/anicine-torrent/internal/analyze"
	"github.com/anicine/anicine-torrent/models"
)

const (
	subs_please_toorent = "subs_please-torrent"
	subs_please_anime   = "subs_please-anime"
	subs_please_path    = "subsplease"
)

func parseSubsPlease(i1 string) (*models.Query, error) {

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
		} else {
			keys = append(keys, v1)
		}
	}

	torrent.Sub = append(torrent.Sub, "ENG")

	i1 = strings.TrimSpace(i1)

	for _, v1 := range keys {
		v1 = strings.ToLower(strings.TrimSpace(v1))
		if strings.Contains(v1, "please") {
			continue
		}
		if len(v1) == 10 {
			torrent.CRC = analyze.CleanSymbols(v1)
			break
		}
	}

	item := new(models.Query)
	item.TorrentFile = torrent

	for _, v1 := range keys {
		v1 = analyze.CleanSymbols(v1)
		v1 = strings.Replace(v1, "-", "~", 1)
		v1 = strings.Replace(v1, "_", "~", 1)
		temp = strings.ReplaceAll(v1, " ", "")

		if ok := analyze.IsValidEpStr(temp); ok {
			eps := analyze.ExtractFloatsWithRanges(temp)
			item.Episodes = eps
			break
		}
	}

	globals = strings.Split(i1, "-")

	for i := len(globals) - 1; i >= 0; i-- {
		if temp = analyze.CleanExtension(globals[i]); temp == "" {
			continue
		}

		globals[i] = temp

		if item.Episodes == nil {
			if globals[i] == globals[0] {
				item.Movie = true
				item.Episodes = []float64{1}
			} else {
				temp = strings.ReplaceAll(globals[i], " ", "")
				temp = strings.ReplaceAll(temp, "_", "~")
				if ok := analyze.IsValidEpStr(temp); ok {
					eps := analyze.ExtractFloatsWithRanges(temp)
					if len(eps) > 1 {
						item.Batch = true
					}
					item.Episodes = eps
					continue
				}
			}
		}

		if i <= len(globals)-1 {
			item.AnimeTitle = new(models.AnimeTitle)
			if len(globals) == 1 {
				item.RomanjiTitle = strings.TrimSpace(globals[i])
				item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
			} else {
				item.RomanjiTitle = strings.Join(globals[:i+1], "-")
				item.RomanjiTitle = strings.TrimSpace(analyze.CleanExtension(item.RomanjiTitle))
			}
			break
		} else {
			return nil, errors.New("cannot parse the anime name")
		}
	}

	return item, nil
}
