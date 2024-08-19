package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestASW(t *testing.T) {

	title := "[ASW] Vivy - Fluorite Eye's Song v3 - S01E10v2 [1080p HEVC x265 10Bit][AAC][0CB1B517]. (batch)"

	item, err := parseASW(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)

}
