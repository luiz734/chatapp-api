package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func saveBytesAsImg(data []byte, filename string) (string, error) {
	var outFileName string
	var outFile *os.File
	var encoderError error

	reader := bytes.NewReader([]byte(data))

	img, format, err := image.Decode(reader)
	if err != nil {
		// log.Printf("Failed to encode and save image: %v", err)
		return "", err
	}

	switch format {
	case "jpeg":
		// s := strconv.Itoa(rand.Int())
		outFileName = fmt.Sprintf("./images/%s.jpg", filename)
		outFile, err = os.Create(outFileName)
		if err != nil {
			// log.Fatalf("Failed to create file: %v", err)
			return "", err
		}
		defer outFile.Close()
		// Set the JPEG compression quality.
		jpegOptions := jpeg.Options{Quality: 80} // Adjust quality as needed.
		encoderError = jpeg.Encode(outFile, img, &jpegOptions)

	case "png":
		// s := strconv.Itoa(rand.Int())
		outFileName = fmt.Sprintf("../chatapp-api/images/%s.png", filename)
		outFile, err = os.Create(outFileName)
		if err != nil {
			// log.Fatalf("Failed to create file: %v", err)
			return "", err
		}
		defer outFile.Close()
		// Set the PNG compression level.
		pngEncoder := png.Encoder{CompressionLevel: png.BestCompression} // Adjust level as needed.
		encoderError = pngEncoder.Encode(outFile, img)

	default:
		return "", errors.New(fmt.Sprintf("Unsupported image format: %v", format))
	}

	if encoderError != nil {
		log.Printf("Failed to encode and save image: %v", encoderError)
		outFileName = ""
		return "", encoderError
	}

	return format, nil
}
