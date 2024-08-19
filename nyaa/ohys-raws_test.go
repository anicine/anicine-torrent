package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOhysRaws(t *testing.T) {

	title := "[ohys-Raws] Attack on Titan: The Final Season - Part 2 - 010.5 (WEB 1080p Hi10 AAC E-AC-3) [Dual-Audio] | Shingeki no Kyojin.mkv"

	item, err := parseOhysRaws(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
