package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJudas(t *testing.T) {

	title := "[Judas] TEST - S1E13.5  | (Test test) [1080p][HEVC x265 10bit][Eng-Subs] (Test Season 1 | S1) (Weekly)"

	item, err := parseJudas(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
