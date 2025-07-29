# go-minimp3

[![Go Reference](https://pkg.go.dev/badge/github.com/cowork-ai/go-minimp3.svg)](https://pkg.go.dev/github.com/cowork-ai/go-minimp3)

[go-minimp3](https://github.com/cowork-ai/go-minimp3) is a Go binding for the [minimp3](https://github.com/lieff/minimp3) C library. The following is the minimp3
description from its author, @lieff.

> Minimalistic, single-header library for decoding MP3. minimp3 is designed to be small, fast (with SSE and NEON support), and accurate (ISO conformant).

go-minimp3 has a very simple interface, one function and one struct, and has zero external dependencies. However, Cgo
must be enabled to compile this package.

## Interface

```go
// Decode parses an MP3 byte slice and returns a [Waveform].
func Decode(mp3Data []byte) (*Waveform, error)

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
func (w *Waveform) NewReader() io.Reader

// WriteTo implements the [io.WriterTo] interface, writing the waveform's samples
// to a writer as signed 16-bit little-endian PCM data.
func (w *Waveform) WriteTo(writer io.Writer) (int64, error)
```

## Examples

### How to convert an MP3 file to a WAV file using [go-audio/wav](https://github.com/go-audio/wav)

Check out [examples/mp3-to-wav](https://github.com/cowork-ai/go-minimp3/blob/main/examples/mp3-to-wav/main.go)

### How to play an MP3 file using [ebitengine/oto](https://github.com/ebitengine/oto)

Check out [examples/play-mp3](https://github.com/cowork-ai/go-minimp3/tree/main/examples/play-mp3/main.go)

## Taskfile.yml

Many useful commands are in two `Taskfile.yml` files: [Taskfile.yml](https://github.com/cowork-ai/go-minimp3/blob/main/Taskfile.yml) and [examples/Taskfile.yml](https://github.com/cowork-ai/go-minimp3/blob/main/examples/Taskfile.yml). To run the tasks, you need to install [go-task/task](https://github.com/go-task/task), which works similarly to [GNU Make](https://www.gnu.org/software/make/).

## Dockerfile

Check out the [Dockerfile](https://github.com/cowork-ai/go-minimp3/blob/main/Dockerfile) for an example of using `golang:1.24` and `gcr.io/distroless/base-debian12` to run `go-minimp3` with Cgo enabled.

```bash
docker build -t cowork-ai/go-minimp3 .
cat ./testdata/44khz128kbps.mp3 | docker run --rm -i cowork-ai/go-minimp3 | ffplay -autoexit -i pipe:
```

Note that the `gcr.io/distroless/static-debian12` image does not work because it lacks `glibc`.
