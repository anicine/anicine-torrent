package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIrizaRaws(t *testing.T) {

	title := "[IrizaRaws] Violet Evergarden - Eien to Jidou Shuki Ningyou (BDRip 1920x1080 x264 10bit FLAC DTSx3)"

	item, err := parseIrizaRaws(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
