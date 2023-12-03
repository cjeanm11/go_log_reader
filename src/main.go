package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	sourcePtr := flag.String("source", "", "Source file")
	outputPtr := flag.String("output", "", "Output file")
	arrayFlag := flag.Bool("array", false, "Output json as a massive array")
	messagePtr := flag.Int("message", 0, "Print a specific message")
	binaryFlag := flag.Bool("binary", false, "Print the raw binary")
	flag.Parse()

	var sourceStream io.Reader = os.Stdin
	var destinationStream io.Writer = os.Stdout

	if *sourcePtr != "" {
		sourceFile, err := os.Open(*sourcePtr)
		if err != nil {
			fmt.Println("Error opening source file:", err)
			os.Exit(1)
		}
		defer sourceFile.Close()
		sourceStream = sourceFile
	}

	if *outputPtr != "" {
		destinationFile, err := os.Create(*outputPtr)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer destinationFile.Close()
		destinationStream = destinationFile
	}

	decoder := json.NewDecoder(sourceStream)
	i := 0

	delimiter := "\n"
	if *arrayFlag {
		fmt.Fprint(destinationStream, "[")
		delimiter = ","
	}

	for {
		var jsonData interface{}
		err := decoder.Decode(&jsonData)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			os.Exit(1)
		}
		i++

		if *messagePtr != 0 {
			if *messagePtr == i {
				break
			}
		}

		if *binaryFlag {
			if i > 1 {
				fmt.Fprintln(destinationStream)
			}
			io.Copy(destinationStream, sourceStream)
		} else {
			jsonDataBytes, err := json.Marshal(jsonData)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				os.Exit(1)
			}

			if i > 1 {
				fmt.Fprint(destinationStream, delimiter)
			}

			fmt.Fprintf(destinationStream, "%s", jsonDataBytes)

			if !*arrayFlag {
				fmt.Fprintln(destinationStream)
			}
		}
	}

	if *binaryFlag {
		fmt.Fprintln(destinationStream)
	} else if *arrayFlag {
		fmt.Fprint(destinationStream, "]")
	} else {
		fmt.Fprintln(destinationStream)
	}
}
