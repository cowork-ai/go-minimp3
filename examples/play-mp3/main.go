package main

import (
	"os"
	"time"

	"github.com/cowork-ai/go-minimp3"
	"github.com/ebitengine/oto/v3"
)

// NOTE: The following code is copied from https://github.com/ebitengine/oto?tab=readme-ov-file#usage.
func main() {
	// Read the mp3 file into memory
	fileBytes, err := os.ReadFile("../testdata/piano.mp3")
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}
	// Decode file
	wav, err := minimp3.Decode(fileBytes)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}
	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.
	op := &oto.NewContextOptions{}
	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = wav.SampleRate
	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = wav.Channels
	// Format of the source. go-minimp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE
	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan
	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(wav.NewReader())
	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()
	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
	// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
	// if err != nil{
	//     panic("player.Seek failed: " + err.Error())
	// }
	// println("Player is now at position:", newPos)
	// player.Play()
	// If you don't want the player/sound anymore simply close
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
