package main

import (
	"encoding/binary"
	"errors"
)

func stm32crc(buf []byte) (crc uint32, err error) {
	crc = 0xFFFFFFFF

	if len(buf)&0x3 != 0 {
		return 0, errors.New("buffer length must be multiple of 4 bytes")
	}

	for i := 0; i < len(buf); i += 4 {
		crc ^= binary.LittleEndian.Uint32(buf[i:])
		for x := 0; x < 32; x++ {
			if crc&0x80000000 != 0 {
				crc = crc << 1 ^ 0x04c11db7
			} else {
				crc = crc << 1
			}
		}
	}
	return
}
