// Command-line Go script to external-sort a large text file.
// Usage:
//
//	go run sortfile.go input.txt output.txt [chunkSizeBytes]
//
// Example:
//
//	go run sortfile.go coupons_unsorted.txt coupons_sorted.txt 100000000
//
// Default chunk size: 100MB
package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run sortfile.go <input> <output> [chunkSizeBytes]")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	chunkSize := 100 * 1024 * 1024 // default 100MB
	if len(os.Args) > 3 {
		if val, err := strconv.Atoi(os.Args[3]); err == nil {
			chunkSize = val
		}
	}

	fmt.Printf("Sorting file: %s → %s (chunk size: %d bytes)\n", inputPath, outputPath, chunkSize)

	if err := ExternalSort(inputPath, outputPath, chunkSize); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	inputLines, err := countLines(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to count input lines: %v\n", err)
		os.Exit(1)
	}
	outputLines, err := countLines(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to count output lines: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Validation: input lines = %d, output lines = %d\n", inputLines, outputLines)
	if inputLines != outputLines {
		fmt.Fprintf(os.Stderr, "Line count mismatch! Possible data loss during sort.\n")
		os.Exit(1)
	}

	fmt.Println("Sorting complete and validated successfully!")
}

func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

// ExternalSort splits, sorts, and merges a large file without loading it fully into memory.
func ExternalSort(inputPath, outputPath string, chunkSize int) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	var tempFiles []string
	scanner := bufio.NewScanner(input)
	var lines []string
	currentSize := 0

	// Split into sorted chunks
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		currentSize += len(line)
		if currentSize >= chunkSize {
			tmp, err := writeSortedChunk(lines)
			if err != nil {
				return err
			}
			tempFiles = append(tempFiles, tmp)
			lines = nil
			currentSize = 0
		}
	}
	if len(lines) > 0 {
		tmp, err := writeSortedChunk(lines)
		if err != nil {
			return err
		}
		tempFiles = append(tempFiles, tmp)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err = mergeSortedFiles(tempFiles, outputPath)

	for _, f := range tempFiles {
		os.Remove(f)
	}
	return err
}

func writeSortedChunk(lines []string) (string, error) {
	sort.Strings(lines)
	tmp, err := os.CreateTemp("", "chunk-*.txt")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	w := bufio.NewWriter(tmp)
	for _, l := range lines {
		_, _ = w.WriteString(l + "\n")
	}
	w.Flush()
	return tmp.Name(), nil
}

type fileScanner struct {
	scanner *bufio.Scanner
	file    *os.File
	value   string
}

type fileHeap []*fileScanner

func (h fileHeap) Len() int           { return len(h) }
func (h fileHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h fileHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *fileHeap) Push(x interface{}) {
	*h = append(*h, x.(*fileScanner))
}
func (h *fileHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func mergeSortedFiles(files []string, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()
	writer := bufio.NewWriter(out)

	h := &fileHeap{}
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		sc := bufio.NewScanner(f)
		if sc.Scan() {
			fs := &fileScanner{scanner: sc, file: f, value: sc.Text()}
			heap.Push(h, fs)
		}
	}

	defer func() {
		for _, fs := range *h {
			fs.file.Close()
		}
	}()

	for h.Len() > 0 {
		fs := heap.Pop(h).(*fileScanner)
		_, _ = writer.WriteString(fs.value + "\n")

		if fs.scanner.Scan() {
			fs.value = fs.scanner.Text()
			heap.Push(h, fs)
		} else {
			fs.file.Close()
		}
	}

	return writer.Flush()
}
