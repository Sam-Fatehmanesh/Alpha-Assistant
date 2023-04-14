package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func float32ToInt8Bytes(values []float32) ([]byte, error) {
	// Create a new slice of int8 values with the same length as the input slice
	int8Values := make([]int8, len(values))

	// Convert each float32 value into an int8 value
	for i, v := range values {
		if v > math.MaxInt8 || v < math.MinInt8 {
			return nil, fmt.Errorf("value at index %d is out of range for int8", i)
		}
		int8Values[i] = int8(v)
	}

	// Encode the int8 slice as a byte array
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, int8Values)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func lastNFalse(b []bool, n int) bool {
	if n > len(b) {
		return false
	}
	for i := len(b) - 1; i >= len(b)-n; i-- {
		if b[i] {
			return false
		}
	}
	return true
}
