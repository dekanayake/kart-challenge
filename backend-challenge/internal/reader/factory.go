package reader

import "fmt"

const (
	HDDReader = "hdd"
	SSDReader = "ssd"
)

func GetFileReader(readerType, rootPath string, chunkSize, searchBatch int) (FileReader, error) {
	switch readerType {
	case HDDReader:
		hddReader, err := newHDDFileReader(rootPath, chunkSize, searchBatch)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize HDDFileReader: %w", err)
		}
		return hddReader, nil

	case SSDReader:
		return nil, fmt.Errorf("Not implemented")

	default:
		return nil, fmt.Errorf("unsupported reader type: %s", readerType)
	}
}
