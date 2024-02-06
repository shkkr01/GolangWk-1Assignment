package main

import (
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

var key string
var inputFile string
var outputFile string
var encryptMode bool
var chunkSize int

// ProcessedChunk holds the processed data along with its index
type ProcessedChunk struct {
	Index int
	Data  []byte
}

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
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	var mu sync.Mutex // Mutex to protect counter
	var counter int

	var wg sync.WaitGroup
	var processedChunks []ProcessedChunk

	for i := 0; ; i++ {
		buffer := make([]byte, chunkSize)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		chunk := buffer[:n]

		wg.Add(1)
		go func(index int, data []byte) {
			defer wg.Done()

			var processedChunk []byte
			var err error
			if encryptMode {
				processedChunk, err = performEncrypt(data, []byte(key))
			} else {
				processedChunk, err = performDecrypt(data, []byte(key))
			}

			if err != nil {
				fmt.Println("Error processing chunk:", err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			// Append the processed chunk with its index
			processedChunks = append(processedChunks, ProcessedChunk{Index: index, Data: processedChunk})
			counter++
		}(i, chunk)
	}

	wg.Wait() // Wait for all processing goroutines to finish

	// Sort processed chunks based on their index
	sort.Slice(processedChunks, func(i, j int) bool {
		return processedChunks[i].Index < processedChunks[j].Index
	})

	// Write the sorted processed chunks to the output file
	for _, pc := range processedChunks {
		_, err := output.Write(pc.Data)
		if err != nil {
			fmt.Println("Error writing processed chunk to file:", err)
			return err
		}
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
