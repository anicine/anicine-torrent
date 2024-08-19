package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMoozzi2(t *testing.T) {

	title := "[Moozzi2] Sokushi Cheat ga Saikyou Sugite BD-BOX (BD 1920x1080 HEVC-YUV444P10 FLACx2) - TV + SP + 4K"

	item, err := parseMoozzi2(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
