package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYameii(t *testing.T) {

	title := "[Yameii] TEST - S1E13.5  | Test test [English Dub] [1080p] [102AE8C4] (Test Season 1 | S1)"

	item, err := parseYameii(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)

}
