package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/h2non/bimg"
)

// Dump :nodoc:
func Dump(data interface{}) string {
	dataByte, _ := json.Marshal(data)
	return string(dataByte)
}

func ParseID(n string) int64 {
	id, _ := strconv.ParseInt(n, 10, 64)
	return id
}

func CompressImage(buffer []byte, quality int) (io.Reader, error) {
	compressed, err := bimg.NewImage(buffer).Process(bimg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(compressed), nil
}
