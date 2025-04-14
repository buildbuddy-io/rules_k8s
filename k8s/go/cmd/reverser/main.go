package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	// 1. Read all data from standard input
	inputData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading standard input: %v", err)
	}

	// 2. Decode all YAML documents from the input
	decoder := yaml.NewDecoder(bytes.NewReader(inputData))
	var documents []interface{} // Use interface{} to hold arbitrary YAML structure

	for {
		var doc interface{}
		err := decoder.Decode(&doc)
		if err == io.EOF {
			// End of input stream
			break
		}
		if err != nil {
			log.Fatalf("Error decoding YAML document: %v", err)
		}
		// Add the successfully decoded document to our slice
		documents = append(documents, doc)
	}

	// 3. Reverse the order of the documents in the slice
	// In-place reversal
	for i, j := 0, len(documents)-1; i < j; i, j = i+1, j-1 {
		documents[i], documents[j] = documents[j], documents[i]
	}

	// 4. Encode the reversed documents back to YAML and write to standard output
	encoder := yaml.NewEncoder(os.Stdout)
	// Configure the encoder for standard YAML formatting (optional, defaults are often fine)
	encoder.SetIndent(2) // Example: Set indentation to 2 spaces

	for i, doc := range documents {
		err := encoder.Encode(doc)
		if err != nil {
			log.Fatalf("Error encoding YAML document #%d: %v", i+1, err)
		}
		// Note: The yaml.v3 Encoder automatically handles the '---' separator
		// between documents when Encode is called multiple times.
	}

	// 5. Ensure everything is written (especially if stdout is buffered elsewhere)
	// Although Encode usually flushes, explicitly closing can catch final errors.
	err = encoder.Close()
	if err != nil {
		// Log error, but don't necessarily fatalf as output might be partially written
		fmt.Fprintf(os.Stderr, "Error closing YAML encoder: %v\n", err)
		// os.Exit(1) // Optionally exit with error code if strictness is required
	}
}
