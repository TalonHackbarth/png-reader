package main

import (
	"fmt"
	"os"
	"path/filepath"
	"png-reader/png"
	"strings"
)

func main() {
	runSingle("basi0g01.png")
	// runBatch("b")
}

func runBatch(prefix string) {
	entries, err := os.ReadDir("test/imageSuite")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), prefix) {
			continue
		}
		file, err := os.ReadFile(filepath.Join("test/imageSuite", entry.Name()))
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		_, err = png.ReadSignature([8]byte(file[:8]))
		if err != nil {
			panic(err)
		}

		chunks, err := png.ReadChunks(file)
		if err != nil {
			panic(err)
		}

		for _, c := range chunks {
			fmt.Println(c)
		}
	}
}

func runSingle(imageName string) {
	file, err := os.ReadFile(filepath.Join("test/imageSuite/", imageName))
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	_, err = png.ReadSignature([8]byte(file[:8]))
	if err != nil {
		panic(err)
	}

	chunks, err := png.ReadChunks(file)
	if err != nil {
		panic(err)
	}

	for _, c := range chunks {
		fmt.Println(c)
	}
}
