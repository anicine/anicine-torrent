package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLostYears(t *testing.T) {

	title := "[LostYears] Attack on Titan: The Final Season - Part 2 (WEB 1080p Hi10 AAC E-AC-3) [Dual-Audio] | Shingeki no Kyojin"

	item, err := parseLostYears(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
