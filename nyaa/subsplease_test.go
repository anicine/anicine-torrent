package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubsPlease(t *testing.T) {

	title := "[SubsPlease] Jujutsu Kaisen S1 - 1~10 (1080p) [Batch] [E81705D2].mkv"

	item, err := parseSubsPlease(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)

}
