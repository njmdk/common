package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPHeader(t *testing.T) {
	r := require.New(t)
	testData := " 1"
	r.Equal("1", TrimLeftRight(testData))
	testData = " 1 \n\t"
	r.Equal("1", TrimLeftRight(testData))
	testData = "  \n  1 1  \t  \n "
	r.Equal("1 1", TrimLeftRight(testData))
}

func BenchmarkNewHTTPHeader(b *testing.B) {
	testData := "  \n  1 1  \t  \n "
	for i := 0; i < b.N; i++ {
		TrimLeftRight(testData)
	}
}
