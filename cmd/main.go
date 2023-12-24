package main

import (
	"fmt"
	gopng "image/png"
	"os"
	"path/filepath"
	"png-reader/image"
	"png-reader/png"
	"strings"
)

func main() {
	// runSingle("basi0g01.png")
	// runBatch("b")

	x := image.Image{Data: [][]uint8{
		{1, 1, 1, 2, 2, 2, 0, 255, 255},
		{0, 0, 50, 0, 0, 20, 0, 0, 50},
		{0, 0, 255, 0, 0, 255, 0, 0, 255},
		{0, 0, 255, 0, 0, 255, 0, 0, 255},
		{0, 0, 255, 0, 0, 255, 0, 0, 255},
		{0, 0, 255, 0, 0, 255, 0, 0, 255},
		{0, 0, 255, 0, 0, 255, 0, 0, 255},
		{0, 0, 255, 0, 0, 255, 0, 0, 255},
		{0, 0, 255, 0, 0, 255, 0, 0, 255}}, Channels: 1}
	y := x.Convert()
	fmt.Println(y)
	f, _ := os.Create("output.png")
	err := gopng.Encode(f, y)
	if err != nil {
		panic(err)
	}
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
			if c.TypeCode == [4]byte{73, 72, 68, 82} {
				id, err := png.ParseIHDR(c)
				if err != nil {
					panic(err)
				}
				fmt.Println(id)
			}
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
		// fmt.Println(c)
		if c.TypeCode == [4]byte{73, 72, 68, 82} {
			id, err := png.ParseIHDR(c)
			if err != nil {
				panic(err)
			}
			fmt.Println(id)
		}
	}
}
