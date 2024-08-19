package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReinForce(t *testing.T) {

	title := "[ReinForce] Uma Musume Pretty Derby S3 Vol 1 (BDRip 1920x1080 x264 FLAC)"

	item, err := parseReinForce(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)
}
