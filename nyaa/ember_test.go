package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmber(t *testing.T) {

	title := "[EMBER] Tokyo Revengers (2021-2023) (Season 1 + 2) (Uncensored) [BDRip] [1080p Dual Audio HEVC 10 bits DDP] (Batch)"

	item, err := parseEmber(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
