package minimp3

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

var writeGolden = flag.Bool("write-golden", false, "")

func TestDecodeMP3(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			"Piano",
			"./testdata/piano.mp3",
			"./testdata/piano.wav",
		},
		{
			"44kHz128kbps",
			"./testdata/44khz128kbps.mp3",
			"./testdata/44khz128kps.wav",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := os.ReadFile(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			pcm, err := DecodeMP3(bs)
			if err != nil {
				t.Fatal(err)
			}
			buf := newIntBuffer(pcm)
			if *writeGolden {
				f, err := os.Create(tt.out)
				if err != nil {
					t.Fatal(err)
				}
				encoder := wav.NewEncoder(f, buf.Format.SampleRate, buf.SourceBitDepth, buf.Format.NumChannels, 1)
				if err := encoder.Write(buf); err != nil {
					t.Fatal(err)
				}
				if err := encoder.Close(); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
			}
			f, err := os.CreateTemp(t.TempDir(), "decode-mp3-*.wav")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())
			encoder := wav.NewEncoder(f, buf.Format.SampleRate, buf.SourceBitDepth, buf.Format.NumChannels, 1)
			if err := encoder.Write(buf); err != nil {
				t.Fatal(err)
			}
			if err := encoder.Close(); err != nil {
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

func newIntBuffer(pcm *PCM) *audio.IntBuffer {
	data := make([]int, len(pcm.Data))
	for i, b := range pcm.Data {
		data[i] = int(b)
	}
	return &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: pcm.NumChannels,
			SampleRate:  pcm.SampleRate,
		},
		Data:           data,
		SourceBitDepth: 16,
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
