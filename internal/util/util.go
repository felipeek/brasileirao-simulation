package util

import (
	"io"
	"math"
	"math/rand"
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

func UtilIntAbs(x int) int {
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

func UtilRandomChoiceStr(choices ...string) interface{} {
	if len(choices) == 0 {
		return nil // Retorna nil se não houver argumentos
	}
	randomIndex := rand.Intn(len(choices))
	return choices[randomIndex]
}

func UtilRandomChoice(choices ...interface{}) interface{} {
	if len(choices) == 0 {
		return nil // Retorna nil se não houver argumentos
	}
	randomIndex := rand.Intn(len(choices))
	return choices[randomIndex]
}
