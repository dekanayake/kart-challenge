package reader

import "fmt"

const (
	HDDReader = "hdd"
	SSDReader = "ssd"
)

func GetFileReader(readerType, rootPath string, chunkSize, searchWorkerPool int) (FileReader, error) {
	switch readerType {
	case HDDReader:
		hddReader, err := newHDDFileReader(rootPath, chunkSize, searchWorkerPool)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize HDDFileReader: %w", err)
		}
		return hddReader, nil

	// Since  SSD  can randoly access file content with less latency , a bainary search directly on file content  can be implemented on SSDReader
	// With this appraoch , nned to find the last record count , then find the midle record , compare that promoCode , and move the read until mathing record is found in binary
	// search on file.
	case SSDReader:
		return nil, fmt.Errorf("Not implemented")

	default:
		return nil, fmt.Errorf("unsupported reader type: %s", readerType)
	}
}
