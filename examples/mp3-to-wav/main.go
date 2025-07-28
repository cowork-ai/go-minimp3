package main

import (
	"io"
	"log"
	"os"

	"github.com/cowork-ai/go-minimp3"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	w, err := minimp3.Decode(bs)
	if err != nil {
		return err
	}
	buf := newIntBuffer(w)
	tmp, err := os.CreateTemp("", "mp3_to_wav_*.wav")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	// wav.NewEncoder requires io.WriterSeeker but os.Stdout is not one even though it pretends like one.
	// so write to a temporary file first.
	encoder := wav.NewEncoder(tmp, buf.Format.SampleRate, buf.SourceBitDepth, buf.Format.NumChannels, 1)
	if err := encoder.Write(buf); err != nil {
		return err
	}
	if err := encoder.Close(); err != nil {
		return err
	}
	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, tmp)
	return err
}

func newIntBuffer(w *minimp3.Waveform) *audio.IntBuffer {
	data := make([]int, len(w.Samples))
	for i, b := range w.Samples {
		data[i] = int(b)
	}
	return &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: w.Channels,
			SampleRate:  w.SampleRate,
		},
		Data:           data,
		SourceBitDepth: 16,
	}
}
