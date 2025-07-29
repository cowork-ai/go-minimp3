package minimp3

import (
	"errors"
	"flag"
	"os"
	"testing"
)

var writeGolden = flag.Bool("write-golden", false, "")

func TestDecodeNoData(t *testing.T) {
	_, err := Decode(nil)
	if got, want := err, ErrNoData; !errors.Is(got, want) {
		t.Errorf("Decode=%v, want=%v", got, want)
	}
	_, err = Decode([]byte{})
	if got, want := err, ErrNoData; !errors.Is(got, want) {
		t.Errorf("Decode=%v, want=%v", got, want)
	}
}

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
			if got, want := fileSize(t, f.Name()), fileSize(t, tt.out); got != want {
				t.Errorf("fileSize=%v, want=%v", got, want)
			}
		})
	}
}

func fileSize(t *testing.T, filename string) int64 {
	t.Helper()
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	return info.Size()
}
