package nyaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEraiRaws(t *testing.T) {

	title := "[Erai-raws] Test - Season - 1 - 1 ~ 10.5v2 (AAC)[HEVC][1080p][Multiple Subtitle][ENG][POR-BR][SPA-LA][SPA][FRE][GER][ITA][RUS][12345678]"

	item, err := parseEraiRaws(title)
	require.Nil(t, err)
	require.NotNil(t, item)

	fmt.Printf("%+v\n", item)
	fmt.Printf("%+v\n", item.AnimeTitle)
	fmt.Printf("%+v\n", item.TorrentFile)
	fmt.Printf("%+v\n", item.Episodes)

}
