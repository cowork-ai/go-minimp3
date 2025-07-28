package minimp3

// #define MINIMP3_IMPLEMENTATION
// #include "minimp3.h"
// #include "minimp3_ex.h"
//
// int decode_mp3(mp3dec_file_info_t* out, const uint8_t* data, size_t data_size) {
//   mp3dec_t mp3d;
//   return mp3dec_load_buf(&mp3d, data, data_size, out, 0, 0);
// }
import "C"

import (
	"fmt"
	"unsafe"
)

type PCM struct {
	NumChannels int
	SampleRate  int
	Data        []int16
}

func DecodeMP3(in []byte) (*PCM, error) {
	var info C.mp3dec_file_info_t
	defer C.free(unsafe.Pointer(info.buffer))
	if errCode := C.decode_mp3(&info, (*C.uint8_t)(&in[0]), C.size_t(len(in))); errCode != 0 {
		return nil, fmt.Errorf("decode_mp3 failed. errCode: %d", errCode)
	}
	data := make([]int16, info.samples)
	copy(data, unsafe.Slice((*int16)(info.buffer), info.samples))
	return &PCM{int(info.channels), int(info.hz), data}, nil
}
