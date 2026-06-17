package main

import (
	"bytes"
	"compress/zlib"
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
		resultString := string(decompressedBytes)
		fmt.Print(resultString)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
