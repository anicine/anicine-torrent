package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeatriceRaws(t *testing.T) {
	title := `Getsuyoubi no Tawawa 2nd \ Tawawa on Monday 2 [BDRip 1920x1080 HEVC TrueHD]`

	item, err := parseBeatriceRaws(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
