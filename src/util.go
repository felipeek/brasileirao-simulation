package main

import (
	"io"
	"math"
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

func UtilMaxInt64(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func UtilClamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}
