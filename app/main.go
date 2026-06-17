package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func createFilePath(hash string) string {
	//first 2 chars → a subdirectory under .git/objects/  remaining 38 chars → the filename inside it
	path := ".git/objects/" + hash[:2] + "/" + hash[2:]
	return path
}

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":

		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		//command structure is cat-file <type> <hash>
		hash := os.Args[3]
		path := ".git/objects/" + hash[:2] + "/" + hash[2:]
		//read file contents
		uncompressedContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(os.Stderr, "Error reading file: %s\n", err)
			os.Exit(1)
		}
		//decompress file with zlib
		reader, err := zlib.NewReader(bytes.NewReader(uncompressedContent))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decompressing file: %s\n", err)
			os.Exit(1)
		}
		//defer will close the reader before the function returns
		defer reader.Close()
		decompressedBytes, err := io.ReadAll(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decompressing file: %s\n", err)
			os.Exit(1)
		}

		//format blob <size>\0<content>
		//find null byte
		nullByteIndex := bytes.IndexByte(decompressedBytes, 0)
		if nullByteIndex == -1 {
			fmt.Fprintf(os.Stderr, "Error finding null byte in decompressed file\n")
			os.Exit(1)
		}
		//content is everything after the null byte
		resultString := string(decompressedBytes[nullByteIndex+1:])
		fmt.Print(resultString)
	case "hash-object":
		commandOne := os.Args[2]

		if commandOne == "-w" && len(os.Args) < 4 {
			fmt.Println("Please provide a file path")
			os.Exit(1)
		}
		var filePath string
		if commandOne == "-w" {
			filePath = os.Args[3]
		} else {
			filePath = os.Args[2]
		}
		//the object file is stored with zlib compression,
		//SHA-1 hash needs to be computed over the "uncompressed" contents of the file, not the compressed version.
		//input for the SHA-1 hash is the header (blob <size>\0) + the actual contents of the file,
		//blob <size>\0<content>

		//build the hash
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			os.Exit(1)
		}
		header := fmt.Sprintf("blob %d\x00", len(fileBytes))
		hashInput := append([]byte(header), fileBytes...)
		hash := sha1.Sum(hashInput)

		if commandOne == "-w" {
			//compress file with zlib
			writer := zlib.NewWriter(os.Stdout)
			defer writer.Close()
			writer.Write(hashInput)
			//write to .git/objects/first2/remaining38
			storagePath := createFilePath(fmt.Sprintf("%x", hash))
			if err := os.WriteFile(storagePath, hashInput, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
			}

		}

		fmt.Printf("%x\n", hash)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
