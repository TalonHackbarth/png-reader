package test

import (
	"errors"
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
