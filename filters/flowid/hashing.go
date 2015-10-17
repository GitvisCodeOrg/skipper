package flowid

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
)

const (
	maxLength       = 64
	minLength       = 8
	flowIdAlphabet  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+"
	alphabetBitMask = 63
)

var (
	ErrInvalidLen = errors.New(fmt.Sprintf("Invalid length. Must be between %d and %d", minLength, maxLength))
	flowIdRegex   = regexp.MustCompile(`^[0-9a-zA-Z+-]+$`)
)

// newFlowId returns a random flowId using the flowIdAlphabet with length l
// The alphabet is limited to 64 elements and requires a random 6 bit value to index any of them
// The cost to rnd.IntXX is not very relevant but the bit shifting operations are faster
// For this reason a single call to rnd.Int63 is used and its bits are mapped up to 8 chunks of 8 bits each
// A final right shift truncates the extra unnecessary bits
func newFlowId(l int) (string, error) {
	if l < minLength || l > maxLength {
		return "", ErrInvalidLen
	}

	u := make([]byte, l)
	for i := 0; i < l; i += 10 {
		b := rand.Int63()
		for e := 0; e < 10 && i+e < l; e++ {
			c := byte(b>>uint(6*e)) & alphabetBitMask // 6 bits only
			u[i+e] = flowIdAlphabet[c]
		}
	}

	return string(u), nil
}

func isValid(flowId string) bool {
	return len(flowId) >= minLength && len(flowId) <= maxLength && flowIdRegex.MatchString(flowId)
}
