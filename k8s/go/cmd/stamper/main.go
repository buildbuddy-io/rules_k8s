package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// stringSliceFlag is a custom flag type to handle multiple occurrences
// of the same flag name, accumulating the values in a slice.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	formatFlag     = flag.String("format", "", "The format string containing stamp variables (e.g., {BUILD_USER}).")
	outputFlag     = flag.String("output", "", "The filename into which we write the result.")
	stampInfoFiles stringSliceFlag // Custom flag type for multiple files
)

func init() {
	// Register the custom flag type. The name matches the Python script's argument.
	flag.Var(&stampInfoFiles, "stamp-info-file", "A file from which to read substitutions (key value pairs, space-separated). May be specified multiple times.")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --format=\"Built by {BUILD_USER} on {BUILD_HOST}\" --stamp-info-file=bazel-out/stable-status.txt --stamp-info-file=bazel-out/volatile-status.txt --output=out.txt\n", os.Args[0])

	}
	flag.Parse()

	// --- Validate required flags ---
	if *formatFlag == "" {
		log.Fatal("--format flag is required")
	}
	if *outputFlag == "" {
		log.Fatal("--output flag is required")
	}
	// --- Read stamp variable files ---
	formatArgs := make(map[string]string)
	for _, infofilePath := range stampInfoFiles {
		file, err := os.Open(infofilePath)
		if err != nil {
			log.Fatalf("Error opening stamp file %q: %v", infofilePath, err)
		}
		// Use defer in a closure to capture the correct file variable
		// if we were to put the file processing in a separate function.
		// Here, it's simple enough inline.
		func() {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			lineNumber := 0
			for scanner.Scan() {
				lineNumber++
				line := scanner.Text()
				if line == "" {
					continue // Skip empty lines
				}

				// SplitN ensures we only split on the first space
				parts := strings.SplitN(line, " ", 2)
				if len(parts) != 2 {
					// Use log instead of raising exception for malformed lines, but continue processing
					log.Printf("Warning: Malformed line %d in %q: %q", lineNumber, infofilePath, line)
					continue
				}

				key, value := parts[0], parts[1]
				if _, exists := formatArgs[key]; exists {
					fmt.Fprintf(os.Stderr, "WARNING: Duplicate value for key %q: using %q (from %s)\n", key, value, infofilePath)
				}
				formatArgs[key] = value
			}
			if err := scanner.Err(); err != nil {
				log.Fatalf("Error reading stamp file %q: %v", infofilePath, err)
			}
		}() // Immediately invoke the closure
	}

	// --- Perform substitutions ---
	// Build the list of replacements for strings.NewReplacer
	// It expects ["old1", "new1", "old2", "new2", ...]
	var replacements []string
	for key, value := range formatArgs {
		placeholder := "{" + key + "}"
		replacements = append(replacements, placeholder, value)
	}

	replacer := strings.NewReplacer(replacements...)
	outputContent := replacer.Replace(*formatFlag)

	// --- Write output file ---
	err := os.WriteFile(*outputFlag, []byte(outputContent), 0644) // Use standard permissions
	if err != nil {
		log.Fatalf("Error writing output file %q: %v", *outputFlag, err)
	}
}
