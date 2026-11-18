package rls

import (
	"github.com/autobrr/rls/taginfo"
	"testing"
)

func TestAnimeReleaseSeriesParsing(t *testing.T) {
	p := NewTagParser(taginfo.All(groupInfos()), DefaultLexers()...)
	tests := []struct {
		name    string
		src     string
		series  int
		episode int
	}{
		{
			name:    "erai ordinal season",
			src:     "[Erai-raws] Hitori no Shita: The Outcast 6th Season - 26 (CA) [1080p CR WEBRip HEVC AAC][MultiSub][46D5B0CC]",
			series:  6,
			episode: 26,
		},
		{
			name:    "erai absolute episode remains seasonless",
			src:     "[Erai-raws] Otaku ni Yasashii Gal wa Inai - 12 [1080p CR WEB-DL AVC AAC][MultiSub][84D6D057]",
			series:  0,
			episode: 12,
		},
		{
			name:    "subsplease explicit season",
			src:     "[SubsPlease] Hitori no Shita - The Outcast S6 - 26 (1080p) [8BE4ED09].mkv",
			series:  6,
			episode: 26,
		},
		{
			name:    "toonshub explicit season episode",
			src:     "[ToonsHub] The Outcast S06E26 1080p CR WEB-DL AAC2.0 H.264 (Yi Ren Zhi Xia, Multi-Subs)",
			series:  6,
			episode: 26,
		},
		{
			name:    "long running standalone episode remains seasonless",
			src:     "[SubsPlease] One Piece - 1125 (1080p) [7E631F90].mkv",
			series:  0,
			episode: 1125,
		},
		{
			name:    "erai roman sequel season",
			src:     "[Erai-raws] Clevatess II - 01 [1080p CR WEB-DL AVC AAC][MultiSub][89A0BAFA]",
			series:  2,
			episode: 1,
		},
		{
			name:    "erai roman sequel season with title words",
			src:     "[Erai-raws] Youjo Senki II - 01 [1080p CR WEB-DL AVC AAC][MultiSub][E4363B39]",
			series:  2,
			episode: 1,
		},
		{
			name:    "erai part sequel season",
			src:     "[Erai-raws] Yoroi-Shinden Samurai Troopers Part 2 - 01 [1080p CR WEB-DL AVC AAC][MultiSub][8B20E329]",
			series:  2,
			episode: 1,
		},
		{
			name:    "erai roman season marker before subtitle",
			src:     "[Erai-raws] Mushoku Tensei III: Isekai Ittara Honki Dasu - 02 [1080p CR WEBRip HEVC AAC][MultiSub][DEC19C44]",
			series:  3,
			episode: 2,
		},
		{
			name:    "erai malformed leading bracket remains seasonless",
			src:     "Erai-raws] Yomi no Tsugai - 13 [1080p CR WEBRip HEVC AAC][MultiSub][7C81BBFA]",
			series:  0,
			episode: 13,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := p.ParseRelease([]byte(test.src))
			if r.Type != Episode {
				t.Fatalf("expected type %s, got %s", Episode, r.Type)
			}
			if r.Series != test.series || r.Episode != test.episode {
				t.Fatalf("expected S%02dE%02d, got S%02dE%02d", test.series, test.episode, r.Series, r.Episode)
			}
		})
	}
}

func TestStandaloneEpisodeMarkers(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		episode int
	}{
		{"One Piece EP1160", "One Piece", 1160},
		{"One Piece EP1162", "One Piece", 1162},
		{"OXXX EP2212", "OXXX", 2212},
		{"One Piece ep1163", "One Piece", 1163},
		{"Another Show Episode1164", "Another Show", 1164},
		{"Another Show EP10000", "Another Show", 10000},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ParseString(test.name)
			if r.Type != Episode || r.Title != test.title || r.Series != 0 || r.Episode != test.episode {
				t.Fatalf("got type=%s title=%q series=%d episode=%d; want type=episode title=%q series=0 episode=%d",
					r.Type, r.Title, r.Series, r.Episode, test.title, test.episode)
			}
		})
	}
}

func TestAnimeReleaseSeasonSuffixTitle(t *testing.T) {
	p := NewTagParser(taginfo.All(groupInfos()), DefaultLexers()...)
	tests := []struct {
		src   string
		title string
	}{
		{
			src:   "[Erai-raws] Clevatess II - 01 [1080p CR WEB-DL AVC AAC][MultiSub][89A0BAFA]",
			title: "Clevatess",
		},
		{
			src:   "[Erai-raws] Youjo Senki II - 01 [1080p CR WEB-DL AVC AAC][MultiSub][E4363B39]",
			title: "Youjo Senki",
		},
		{
			src:   "[Erai-raws] Yoroi-Shinden Samurai Troopers Part 2 - 01 [1080p CR WEB-DL AVC AAC][MultiSub][8B20E329]",
			title: "Yoroi-Shinden Samurai Troopers",
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := p.ParseRelease([]byte(test.src))
			if r.Title != test.title {
				t.Fatalf("expected title %q, got %q", test.title, r.Title)
			}
		})
	}
}

func TestOnePiece1000Plus(t *testing.T) {
	const name = "[ToonsHub] One Piece EP1160 1080p CR WEB-DL DUAL AAC2.0 H.264 (Dual-Audio, Multi-Subs) "
	r := ParseString(name)
	if r.Type != Episode || r.Title != "One Piece" || r.Series != 0 || r.Episode != 1160 ||
		r.Source != "WEB-DL" || r.Resolution != "1080p" || r.Collection != "CR" ||
		len(r.Codec) != 1 || r.Codec[0] != "H.264" {
		t.Fatalf("unexpected release: %+v", r)
	}
}
