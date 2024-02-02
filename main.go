package main

import (
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var key string
var inputFile string
var outputFile string
var encryptMode bool
var chunkSize int

func init() {
	flag.StringVar(&key, "key", "", "Encryption key")
	flag.StringVar(&inputFile, "file", "", "Path to the input file")
	flag.StringVar(&outputFile, "out", "", "Path to the output file")
	flag.BoolVar(&encryptMode, "encrypt", false, "Encrypt mode")
	flag.IntVar(&chunkSize, "chunk-size", 1024, "Chunk size in bytes")
	flag.Parse()
}

func performEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, make([]byte, aes.BlockSize))
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

func performDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, make([]byte, aes.BlockSize))
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func processFile() error {
	// Read input file
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	var result []byte

	for i := 0; i < len(inputData); i += chunkSize {
		end := i + chunkSize
		if end > len(inputData) {
			end = len(inputData)
		}
		chunk := inputData[i:end]

		var processedChunk []byte
		if encryptMode {
			processedChunk, err = performEncrypt(chunk, []byte(key))
		} else {
			processedChunk, err = performDecrypt(chunk, []byte(key))
		}

		if err != nil {
			fmt.Println("Error processing chunk:", err)
			return err
		}

		result = append(result, processedChunk...)
	}

	// Write result to the output file 0644 is used for the read write permission
	err = os.WriteFile(outputFile, result, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	start := time.Now()
	if key == "" || inputFile == "" || outputFile == "" {
		fmt.Println("Usage: main.go --file <input_file> --key <encryption_key> --out <output_file> [--encrypt] [--chunk-size <size>]")
		return
	}

	err := processFile()
	if err != nil {
		fmt.Println("Error processing file:", err)
		elapsed := time.Since(start)
		log.Printf("File processing took %s", elapsed)
		return
	}

	fmt.Println("File processed successfully.")
	elapsed := time.Since(start)
	log.Printf("File processing took %s", elapsed)
}
