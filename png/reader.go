package png

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"math"
)

const (
	ctGrayscale      = 0
	ctTrueColor      = 2
	ctIndexedColor   = 3
	ctGrayscaleAlpha = 4
	ctTrueColorAlpha = 6
)

const (
	filterNone    = 0
	filterSub     = 1
	filterUp      = 2
	filterAverage = 3
	filterPaeth   = 4
)

const (
	interlaceNone  = 0
	interlaceAdam7 = 1
)

type ImageData struct {
	Width             uint32
	Height            uint32
	BitDepth          uint8
	ColorType         uint8
	CompressionMethod uint8
	FilterMethod      uint8
	InterlaceMethod   uint8
}

type Chunk struct {
	Length   uint32
	TypeCode [4]byte
	Data     []byte
	CRC      uint32
}

var FormatError = errors.New("invalid png format")
var CRCMismatchError = errors.New("cyclic redundancy check found corrupt data")
var IHDRLenError = errors.New("IHDR chunk length invalid")
var IHDRDimensionError = errors.New("invalid image dimensions")

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

func ParseIHDR(chunk Chunk) (ImageData, error) {
	if chunk.Length != 13 {
		return ImageData{}, IHDRLenError
	}
	if len(chunk.Data) != 13 { // Just ensures no possible issue. Same as last check
		return ImageData{}, IHDRLenError
	}
	if chunk.TypeCode != [4]byte{73, 72, 68, 82} {
		return ImageData{}, errors.New("invalid IHDR code")
	}
	dimensionLimit := uint32(math.Pow(2, 31))

	result := ImageData{}

	result.Width = binary.BigEndian.Uint32(chunk.Data[:4])
	result.Height = binary.BigEndian.Uint32(chunk.Data[4:8])
	if result.Width <= 0 || result.Height <= 0 {
		return ImageData{}, IHDRDimensionError
	}
	if result.Width > dimensionLimit || result.Height > dimensionLimit {
		return ImageData{}, IHDRDimensionError
	}
	result.BitDepth = chunk.Data[8]
	result.ColorType = chunk.Data[9]
	if !ValidateColorDepth(result.ColorType, result.BitDepth) {
		return ImageData{}, errors.New("color type & bit depth combination is either not allowed or invalid")
	}

	result.CompressionMethod = chunk.Data[10]
	if result.CompressionMethod != 0 {
		return ImageData{}, errors.New("unsupported compression method. only png v1.2 deflate compression is valid")
	}
	result.FilterMethod = chunk.Data[11]
	if result.FilterMethod != 0 {
		return ImageData{}, errors.New("unsupported filter method. only png v1.2 filters are valid")
	}
	result.InterlaceMethod = chunk.Data[12]
	if result.InterlaceMethod > interlaceAdam7 {
		return ImageData{}, errors.New("unsupported interlace method. only png v1.2 interlacing (adam7) is valid")
	}

	return result, nil
}

func ValidateColorDepth(colorType uint8, bitDepth uint8) bool {
	/* http://www.libpng.org/pub/png/spec/1.2/PNG-Chunks.html
	   Color    Allowed
	   Type    Bit Depths

	   0       1,2,4,8,16
	   2       8,16
	   3       1,2,4,8
	   4       8,16
	   6       8,16
	*/
	switch bitDepth { // Confirm BitDepth & Color Profile is an allowed combination
	case 1:
		if colorType != ctGrayscale && colorType != ctIndexedColor {
			return false
		}
	case 2:
		if colorType != ctGrayscale && colorType != ctIndexedColor {
			return false
		}
	case 4:
		if colorType != ctGrayscale && colorType != ctIndexedColor {
			return false
		}
	case 8:
		// All Color Types can have a bit depth of 8. Always valid
		switch colorType {
		case ctGrayscale:
			return true
		case ctTrueColor:
			return true
		case ctIndexedColor:
			return true
		case ctGrayscaleAlpha:
			return true
		case ctTrueColorAlpha:
			return true
		default:
			return false // This mean colorType is invalid (not 0, 2, 3, 4, or 6)
		}
	case 16:
		switch colorType {
		case ctGrayscale:
			return true
		case ctTrueColor:
			return true
		case ctIndexedColor:
			return false
		case ctGrayscaleAlpha:
			return true
		case ctTrueColorAlpha:
			return true
		default:
			return false // This mean colorType is invalid (not 0, 2, 3, 4, or 6)
		}
	default:
		return false // This means bitDepth is not one of the 5 valid options
	}
	return true
}
