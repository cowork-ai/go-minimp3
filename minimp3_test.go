package minimp3

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"os"
	"testing"
)

var writeGolden = flag.Bool("write-golden", false, "")

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			"Piano",
			"./testdata/piano.mp3",
			"./testdata/piano_s16le_44.1khz_stereo.pcm",
		},
		{
			"44kHz128kbps",
			"./testdata/44khz128kbps.mp3",
			"./testdata/44khz128kbps_s16le_44.1khz_stereo.pcm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := os.ReadFile(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			wav, err := Decode(bs)
			if err != nil {
				t.Fatal(err)
			}
			if *writeGolden {
				f, err := os.Create(tt.out)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := wav.WriteTo(f); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
			}
			f, err := os.CreateTemp(t.TempDir(), "decode_mp3_s16le_44.1khz_stereo_*.pcm")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())
			if _, err := wav.WriteTo(f); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			if got, want := md5Sum(t, f.Name()), md5Sum(t, tt.out); got != want {
				t.Errorf("md5Sum=%v, want=%v", got, want)
			}
		})
	}
}

func md5Sum(t *testing.T, filename string) string {
	t.Helper()
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}
