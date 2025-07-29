package minimp3

// #define MINIMP3_IMPLEMENTATION
// #include "minimp3.h"
// #include "minimp3_ex.h"
//
// int decode(mp3dec_file_info_t* out, const uint8_t* data, size_t data_size) {
//   mp3dec_t mp3d;
//   return mp3dec_load_buf(&mp3d, data, data_size, out, 0, 0);
// }
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"unsafe"
)

// Decode parses an MP3 byte slice and returns a [Waveform].
func Decode(mp3Data []byte) (*Waveform, error) {
	var info C.mp3dec_file_info_t
	defer C.free(unsafe.Pointer(info.buffer))
	if errCode := C.decode(&info, (*C.uint8_t)(&mp3Data[0]), C.size_t(len(mp3Data))); errCode != 0 {
		return nil, fmt.Errorf("decode failed. errCode: %d", errCode)
	}
	samples := make([]int16, info.samples)
	copy(samples, unsafe.Slice((*int16)(info.buffer), info.samples))
	return &Waveform{
		Channels:   int(info.channels),
		SampleRate: int(info.hz),
		Samples:    samples,
	}, nil
}

// Waveform represents decoded PCM audio data. 🎶
type Waveform struct {
	// Channels is the number of audio channels (e.g., 1 for mono, 2 for stereo).
	Channels int
	// SampleRate is the number of samples per second (e.g., 44100 Hz).
	SampleRate int
	// Samples contains the interleaved audio data.
	Samples []int16
}

// NewReader returns a new [io.Reader] that streams the waveform's data
// as signed 16-bit little-endian PCM.
func (w *Waveform) NewReader() io.Reader {
	var buf bytes.Buffer
	if _, err := w.WriteTo(&buf); err != nil {
		log.Panic(err)
	}
	return &buf
}

// WriteTo implements the [io.WriterTo] interface, writing the waveform's samples
// to a writer as signed 16-bit little-endian PCM data.
func (w *Waveform) WriteTo(writer io.Writer) (int64, error) {
	if err := binary.Write(writer, binary.LittleEndian, w.Samples); err != nil {
		return 0, err
	}
	// Each int16 sample is 2 bytes.
	return int64(len(w.Samples) * 2), nil
}
