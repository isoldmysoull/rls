package rls_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/moistari/rls"
	"github.com/stretchr/testify/assert"
)

func ExampleParse() {
	const title = "The_Velvet_Underground-The_Complete_Matrix_Tapes-Reissue_Limited_Edition_Boxset-8LP-2019-NOiR"
	r := rls.ParseString(title)
	fmt.Printf("%q:\n", r)
	fmt.Printf("  type: %s\n", r.Type)
	fmt.Printf("  artist: %s\n", r.Artist)
	fmt.Printf("  title: %s\n", r.Title)
	fmt.Printf("  year: %d\n", r.Year)
	fmt.Printf("  disc: %s\n", r.Disc)
	fmt.Printf("  source: %s\n", r.Source)
	fmt.Printf("  edition: %q\n", r.Edition)
	fmt.Printf("  other: %q\n", r.Other)
	fmt.Printf("  group: %s\n", r.Group)
	// Output:
	// "The_Velvet_Underground-The_Complete_Matrix_Tapes-Reissue_Limited_Edition_Boxset-8LP-2019-NOiR":
	//   type: music
	//   artist: The Velvet Underground
	//   title: The Complete Matrix Tapes
	//   year: 2019
	//   disc: 8x
	//   source: LP
	//   edition: ["Limited.Edition"]
	//   other: ["REISSUE" "BOXSET"]
	//   group: NOiR
}

func TestOnePiece1000Plus(t *testing.T) {
	const title = "[ToonsHub] One Piece EP1160 1080p CR WEB-DL DUAL AAC2.0 H.264 (Dual-Audio, Multi-Subs) "
	r := rls.ParseString(title)
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "One Piece", r.Title)
	assert.Equal(t, "WEB-DL", r.Source)
	assert.Equal(t, "1080p", r.Resolution)
	assert.Equal(t, "CR", r.Collection)
	assert.Len(t, r.Codec, 1)
	if len(r.Codec) == 1 {
		assert.Equal(t, "H.264", r.Codec[0])
	}

	assert.Equal(t, 1160, r.Episode)
	assert.Zero(t, r.Series)
	fmt.Println(string(data))
}
