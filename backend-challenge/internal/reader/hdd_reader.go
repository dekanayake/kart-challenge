package reader

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dekanayake/kart-challenge/backend-challenge/internal/config"
)

type fileIndex struct {
	path         string
	chunkKeys    []string
	chunkOffsets []int64
	firstKey     string
	lastKey      string
	chunkSize    int
}

// HDDFileReader manages multiple indexed coupon files.
type HDDFileReader struct {
	rootPath    string
	chunkSize   int
	searchBatch int
	fileIndexes []*fileIndex // sorted by firstKey
}

// Coupon File reader for  HDD storage.
// Since the files are large ~1GB , its not scalable and performant to load the contents in files to memory
// Files can be readed in streaming , but the file content will need to be readed from the beginning , until the coupon code is matched, which can be Big(N) time in worse case scenario.
// Following optimisations are implemented to make the search faster
//
//	(1) Sort the contents in the files in ascending order , its assume the files will be in sorting order before the application starts. sort script is available at utils/sort_file.gp
//	(2) When the applcation starts reader  will read all the file contents
//	(3) For each defined chunk size  reader will store the line , and offset of that line , this is called as partial index.
//	(4) Also for each file  reader captures the first line and last line of the file
//	(5) When SearchPromoCode function called , reader will filter files which assumes the content inside the file . To do this it check wheter promoCode >= firstline && promoCode <= lastLine
//	(6) For each filterd files reader will do a bianry search on the partial index for that file . when partial index is found , it gives the indication coupon code may be available in the chunk starts from partil index ,
//	    to cover the range more , the program will read the lines from next partial index as well. These entries will load to memory , but since their size = chunk size * 2 . the memory foot print will be small for the records.
//	(7) The program will do a binary search on loaded lines to find a matching couponCode.
//
// . (8) coupon code search on a file operation will handle concurrently , in a batch . This is to avoid  creating go routines when large set of files are available to search.
// .     Once the condition is matched context will be cancelled so any running coupon code search go routine will stop.
func newHDDFileReader(rootPath string, chunkSize, searchBatch int) (*HDDFileReader, error) {
	files, err := filepath.Glob(filepath.Join(rootPath, "couponbase*"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("no coupon files found")
	}

	var indexes []*fileIndex
	for _, path := range files {
		fi, err := buildPartialIndex(path, chunkSize)
		if err != nil {
			return nil, fmt.Errorf("index build failed for %s: %w", path, err)
		}
		indexes = append(indexes, fi)
	}

	config.Logger.Info().
		Msg("Partial indexes are created.")
	for _, fileIndex := range indexes {
		config.Logger.Info().
			Str("path", fileIndex.path).
			Int("index size", int(len(fileIndex.chunkOffsets))).
			Msg("Partial index information")
	}

	// sort by firstKey lexicographically
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i].firstKey < indexes[j].firstKey
	})

	return &HDDFileReader{
		rootPath:    rootPath,
		chunkSize:   chunkSize,
		searchBatch: searchBatch,
		fileIndexes: indexes,
	}, nil
}

// buildPartialIndex builds a simple index for one coupon file.
func buildPartialIndex(path string, chunkSize int) (*fileIndex, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var firstKey, lastKey string
	var offsets []int64
	var keys []string
	var offset int64
	lineCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if firstKey == "" {
			firstKey = line
		}
		lastKey = line
		if lineCount%chunkSize == 0 {
			offsets = append(offsets, offset)
			keys = append(keys, line)
		}
		offset += int64(len(line) + 1)
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &fileIndex{
		path:         path,
		chunkKeys:    keys,
		chunkOffsets: offsets,
		firstKey:     firstKey,
		lastKey:      lastKey,
		chunkSize:    chunkSize,
	}, nil
}

func (r *HDDFileReader) SearchPromo(ctx context.Context, promo string) (bool, error) {
	if len(r.fileIndexes) == 0 {
		return false, errors.New("no indexed files")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan searchResult, len(r.fileIndexes))
	var wg sync.WaitGroup

	for _, fi := range r.fileIndexes {
		// skip files whose range cannot include the promo
		if promo < fi.firstKey || promo > fi.lastKey {
			config.Logger.Debug().
				Str("File path", fi.path).
				Msg("Skipping searching in the file since Promo code is either small or larger than the first and last record of the file")
			continue
		}

		wg.Add(1)
		go func(fi *fileIndex) {
			config.Logger.Info().
				Str("File path", fi.path).
				Msg("Searching the promo index")

			defer wg.Done()
			ok, err := searchPromoInFile(ctx, fi, promo)
			select {
			case <-ctx.Done():
				return
			case results <- searchResult{found: ok, path: fi.path, err: err}:
			}
		}(fi)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	foundCount := 0
	var errs []error

	for res := range results {
		if res.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", res.path, res.err))
			continue
		}
		if res.found {
			foundCount++
			if foundCount >= 2 {
				cancel()
				break
			}
		}
	}

	if len(errs) > 0 {
		return foundCount >= 2, errors.Join(errs...)
	}
	return foundCount >= 2, nil
}

// searchPromoInFile performs in-memory binary search for the target promo.
func searchPromoInFile(ctx context.Context, fi *fileIndex, promo string) (bool, error) {
	file, err := os.Open(fi.path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	idx := sort.Search(len(fi.chunkKeys), func(i int) bool {
		return fi.chunkKeys[i] >= promo
	})
	if idx > 0 {
		idx--
	}
	if idx >= len(fi.chunkOffsets) {
		return false, nil
	}

	_, err = file.Seek(fi.chunkOffsets[idx], 0)
	if err != nil {
		return false, err
	}

	// Read one chunk fully into memory
	scanner := bufio.NewScanner(file)
	var lines []string
	count := 0
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return false, nil
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
		count++
		if count >= fi.chunkSize*2 {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	// binary search within the file
	i := sort.SearchStrings(lines, promo)
	if i < len(lines) && lines[i] == promo {
		return true, nil
	}

	return false, nil
}

type searchResult struct {
	found bool
	path  string
	err   error
}
