package png

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

type Chunk struct {
	Length   uint32
	TypeCode [4]byte
	Data     []byte
	CRC      uint32
}

var FormatError = errors.New("invalid png format")
var CRCMismatchError = errors.New("cyclic redundancy check found corrupt data")

func ReadSignature(sig [8]byte) (bool, error) { // Checks for PNG Sig. \211 P N G \r \n \032 \n
	if sig == [8]byte{137, 80, 78, 71, 13, 10, 26, 10} {
		return true, nil
	} else {
		return false, FormatError
	}
}

func ReadChunks(data []byte) ([]Chunk, error) {
	result := make([]Chunk, 0)
	startLocation := 8
	for {
		if startLocation > len(data) {
			return nil, FormatError
		}
		chunkLength := binary.BigEndian.Uint32(data[startLocation : startLocation+4])
		chunkCode := [4]byte(data[startLocation+4 : startLocation+8])

		chunkData := data[startLocation+8 : startLocation+8+int(chunkLength)]
		chunkCRC := binary.BigEndian.Uint32(data[startLocation+8+int(chunkLength) : startLocation+8+int(chunkLength)+4])

		crcInput := data[startLocation+4 : startLocation+8+int(chunkLength)]

		result = append(result, Chunk{chunkLength, chunkCode, chunkData, chunkCRC})
		_, err := VerifyCRC(crcInput, chunkCRC)
		if err != nil {
			return nil, err
		}

		if chunkCode == [4]byte{73, 69, 78, 68} { // Check for IEND end of file marker
			break
		}

		startLocation += 8 + int(chunkLength) + 4
	}

	return result, nil
}

func VerifyCRC(data []byte, crc uint32) (bool, error) {
	if crc32.ChecksumIEEE(data) != crc {
		return false, CRCMismatchError
	} else {
		return true, nil
	}
}
