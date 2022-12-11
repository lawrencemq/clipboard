package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func printRed(text string) {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	fmt.Println(string(colorRed), text, string(colorReset))
}

func printGreen(text string) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), text, string(colorReset))
}

func printYellow(text string) {
	colorReset := "\033[0m"
	colorYellow := "\033[33m"
	fmt.Println(string(colorYellow), text, string(colorReset))
}

func getClipboardFilename() string {
	return os.TempDir() + "clipboard.txt"
}

func pipeToClipboard() {

	// open output file
	fo, err := os.Create(getClipboardFilename())
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(fo)
	bytesWritten, err := reader.WriteTo(writer)
	if err != nil {
		panic(err)
	}
	printGreen(strconv.FormatInt(bytesWritten, 10) + " bytes written to clipboard.")
}

func peek() {
	clipboardFilename := getClipboardFilename()
	if _, err := os.Stat(clipboardFilename); err == nil {
		file, err := os.Open(clipboardFilename)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = file.Close(); err != nil {
				panic(err)
			}
		}()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			fmt.Println(scanner.Text()) // token in unicode-char
		}
	} else {
		printYellow("There is nothing in the clipboard.")
	}

}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func getUsage() string {
	return `Usage: clipboard <copy|paste|cut|peek> [filename]`
}

func main() {

	if isInputFromPipe() {
		pipeToClipboard()
	} else {
		argsWithoutProg := os.Args[1:]
		if len(argsWithoutProg) < 1 {
			printRed(getUsage())
			return
		}

		switch command := strings.ToLower(argsWithoutProg[0]); command {
		case "copy":
			printGreen("copy")
		case "paste":
			printGreen("paste")
		case "cut":
			printGreen("cut")
		case "peek":
			peek()
		default:
			printRed("Unknown command '" + command + "'. Options are 'copy', 'paste', 'cut', or 'peek'.")
		}

	}

}
