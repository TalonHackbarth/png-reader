package test

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"png-reader/png"
	"strings"
	"testing"
)

// Test Image suite
// http://www.schaik.com/pngsuite/

// 12-23/2023 : 87.5% Test Coverage

func TestFileSignature(t *testing.T) {
	t.Run("Testing Valid Signature", func(t *testing.T) {
		t.Parallel()
		result, err := png.ReadSignature([8]byte{137, 80, 78, 71, 13, 10, 26, 10})
		if !result {
			t.Error("Signature rest read valid sig as invalid")
		}
		if err != nil {
			t.Error("Error was not nil")
		}
	})
	t.Run("Testing Invalid Signatures", func(t *testing.T) {
		t.Parallel()
		result, err := png.ReadSignature([8]byte{0, 78, 1, 43, 56, 4, 26, 10})
		if result {
			t.Error("Invalid sig read as valid")
		}
		if !errors.Is(err, png.FormatError) {
			t.Error("Did not raise png.FormatError")
		}
	})
	t.Run("Testing Empty Signatures", func(t *testing.T) {
		t.Parallel()
		result, err := png.ReadSignature([8]byte{})
		if result {
			t.Error("Empty sig accepted")
		}
		if !errors.Is(err, png.FormatError) {
			t.Error("Empty sig did not raise error")
		}
	})
}

func TestFileSuite(t *testing.T) { // TODO: Replace with better tests
	t.Run("Testing Valid Chunk Files", func(t *testing.T) {
		t.Parallel()
		entries, err := os.ReadDir("imageSuite")
		if err != nil {
			panic(err)
		}
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "x") {
				continue
			}
			file, err := os.ReadFile(filepath.Join("imageSuite", entry.Name()))
			if err != nil {
				panic(err)
			}

			_, err = png.ReadChunks(file)
			if err != nil {
				t.Error(entry.Name(), err.Error())
			}

		}
	})
	// TODO: Test invalid chunk files (x prefix)
	t.Run("Testing Invalid Chunk Files", func(t *testing.T) {
		t.Parallel()
		entries, err := os.ReadDir("imageSuite")
		if err != nil {
			panic(err)
		}
		for _, entry := range entries {
			if !strings.HasPrefix(entry.Name(), "xcs") ||
				!strings.HasPrefix(entry.Name(), "xhd") ||
				!strings.HasPrefix(entry.Name(), "xlf") {
				continue
			}
			file, err := os.ReadFile(filepath.Join("imageSuite", entry.Name()))
			if err != nil {
				panic(err)
			}

			_, err = png.ReadChunks(file)
			if !errors.Is(err, png.CRCMismatchError) {
				t.Error(entry.Name(), err.Error())
			}

		}
	})
}

func TestIHDR(t *testing.T) {
	validData01 := png.Chunk{
		Length:   13,
		TypeCode: [4]byte{73, 72, 68, 82},
		Data:     []byte{0, 0, 0, 32, 0, 0, 0, 32, 1, 0, 0, 0, 1},
		CRC:      738621391,
	}
	invalidData01 := png.Chunk{
		Length:   12,
		TypeCode: [4]byte{73, 72, 68, 82},
		Data:     []byte{0, 0, 0, 32, 0, 0, 0, 32, 1, 0, 0, 0, 1},
		CRC:      738621391,
	}
	invalidData02 := png.Chunk{
		Length:   13,
		TypeCode: [4]byte{73, 72, 68, 82},
		Data:     []byte{0, 0, 0, 32, 0, 0, 0, 32, 1, 0, 0, 0, 1, 3},
		CRC:      738621391,
	}
	invalidDimension01 := png.Chunk{
		Length:   13,
		TypeCode: [4]byte{73, 72, 68, 82},
		Data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		CRC:      738621391,
	}
	invalidDimension02 := png.Chunk{
		Length:   13,
		TypeCode: [4]byte{73, 72, 68, 82},
		Data:     []byte{0xff, 0xff, 0xff, 0xff, 0, 0, 0, 32, 1, 0, 0, 0, 1},
		CRC:      738621391,
	}
	invalidCode := png.Chunk{
		Length:   13,
		TypeCode: [4]byte{73, 72, 68, 0},
		Data:     []byte{0, 0, 0, 32, 0, 0, 0, 32, 1, 0, 0, 0, 1},
		CRC:      738621391,
	}
	t.Run("Testing Invalid IHDR Data Length", func(t *testing.T) {
		t.Parallel()
		_, err := png.ParseIHDR(invalidData01)
		if !errors.Is(err, png.IHDRLenError) {
			t.Error("Did not catch invalid chunk length", err)
		}
		_, err = png.ParseIHDR(invalidData02)
		if !errors.Is(err, png.IHDRLenError) {
			t.Error("Did not catch invalid chunk length", err)
		}
	})
	t.Run("Testing Valid IHDR Values", func(t *testing.T) {
		t.Parallel()
		_, err := png.ParseIHDR(validData01)
		if err != nil {
			t.Error(err.Error())
		}
	})
	t.Run("Testing Invalid IHDR Dimensions", func(t *testing.T) {
		t.Parallel()
		_, err := png.ParseIHDR(invalidDimension01)
		if !errors.Is(err, png.IHDRDimensionError) {
			t.Error(err.Error())
		}
		_, err = png.ParseIHDR(invalidDimension02)
		fmt.Println(err)
		if !errors.Is(err, png.IHDRDimensionError) {
			t.Error(err.Error())
		}
	})
	t.Run("Testing Invalid IHDR Code", func(t *testing.T) {
		t.Parallel()
		_, err := png.ParseIHDR(invalidCode)
		if err.Error() != "invalid IHDR code" {
			t.Error("Did not catch invalid chunk code", err)
		}
	})
}

func TestColorDepthValidation(t *testing.T) {
	t.Run("Testing Valid Color Type & Bit Depth Values", func(t *testing.T) {
		t.Parallel()
		cases := map[uint8][]uint8{
			0: {1, 2, 4, 8, 16},
			2: {8, 16},
			3: {1, 2, 4, 8},
			4: {8, 16},
			6: {8, 16},
		}
		for colorType, validBitDepths := range cases {
			for _, bd := range validBitDepths {
				if !png.ValidateColorDepth(colorType, bd) {
					t.Error("Color Depth Validation Failed! (", colorType, bd, ")")
				}
			}
		}

	})
	t.Run("Testing Invalid Color Type & Bit Depth Values", func(t *testing.T) {
		t.Parallel()
		cases := map[uint8][]uint8{
			0:  {5, 99, 10, uint8(math.Pow(8, 26)), 0, 255},
			2:  {5, 99, 10, uint8(math.Pow(8, 26)), 0, 255, 2, 4},
			3:  {5, 99, 10, uint8(math.Pow(8, 26)), 0, 255, 16},
			4:  {5, 99, 10, uint8(math.Pow(8, 26)), 0, 255, 2, 4},
			6:  {5, 99, 10, uint8(math.Pow(8, 26)), 0, 255, 2, 4},
			7:  {5, 99, 10, uint8(math.Pow(8, 26)), 0, 255, 16, 4},
			99: {1, 0, 255, 16, 2, 4, 8},
		}
		for colorType, validBitDepths := range cases {
			for _, bd := range validBitDepths {
				if png.ValidateColorDepth(colorType, bd) {
					t.Error("Color Depth Validation Failed! (", colorType, bd, ") FALSE POSITIVE")
				}
			}
		}

	})
}
