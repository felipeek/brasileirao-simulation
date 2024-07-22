package main

import (
	"io"
	"os"
)

func UtilReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}

func UtilInt64Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func MaxInt64(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
