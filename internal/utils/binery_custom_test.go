package utils

import (
	"fmt"
	"testing"
)

func TestEncodeAllStructs(t *testing.T) {
	structs := []struct {
		name  string
		frame Action
	}{
		// {"ReversedFrame", &ReversedFrame{}},
		// {"ShadowFrame", &ShadowFrame{}},
		// {"StealthFrame", &StealthFrame{}},
		// {"XDRFrame", &XDRFrame{}},
		// {"ShuffledFrame", &ShuffledFrame{}},
		// {"XORMaskedFrame", &XORMaskedFrame{}},
		{"BinaryFrame", &BinaryFrame{}},
		{"NormalFrame", &NormalFrame{}},
	}
	testMsg := "Hello, World!"
	testTransport := byte(0x01)

	for _, s := range structs {
		t.Run(fmt.Sprintf("Testing %s", s.name), func(t *testing.T) {
			encoded, err := s.frame.Encode(testMsg, testTransport)
			if err != nil {
				t.Fatalf("Encode failed for %s: %v", s.name, err)
			}

			k, decodedTransport, err := s.frame.Decode(encoded)

			if err != nil {
				t.Fatalf("Decode failed for %s: %v", s.name, err)
			}

			if k != testMsg {
				t.Errorf("Decoded message mismatch for %s. Got: %s, Expected: %s", s.name, k, testMsg)
			}
			if decodedTransport != testTransport {
				t.Errorf("Decoded transport mismatch for %s. Got: %d, Expected: %d", s.name, decodedTransport, testTransport)
			}
		})
	}
}

// func TestEncodeString(t *testing.T) {
// 	structs := []struct {
// 		name  string
// 		frame Action
// 	}{
// 		// {"ReversedFrame", &ReversedFrame{}},
// 		// {"ShadowFrame", &ShadowFrame{}},
// 		// {"StealthFrame", &StealthFrame{}},
// 		// {"XDRFrame", &XDRFrame{}},
// 		// {"ShuffledFrame", &ShuffledFrame{}},
// 		// {"XORMaskedFrame", &XORMaskedFrame{}},
// 		{"BinaryFrame", &BinaryFrame{}},
// 		{"NormalFrame", &NormalFrame{}},
// 	}
// 	testMsg := "Hello, World!"

// 	for _, s := range structs {
// 		t.Run(fmt.Sprintf("Testing %s", s.name), func(t *testing.T) {
// 			encoded, err := s.frame.EncodeString(testMsg)
// 			if err != nil {
// 				t.Fatalf("Encode failed for %s: %v", s.name, err)
// 			}
// 			decodedMsg, _, err := s.frame.DecodeString(encoded)
// 			if err != nil {
// 				t.Fatalf("Decode failed for %s: %v", s.name, err)
// 			}

// 			if decodedMsg != testMsg {
// 				t.Errorf("Decoded message mismatch for %s. Got: %s, Expected: %s", s.name, decodedMsg, testMsg)
// 			}
// 		})
// 	}
// }

// func BenchmarkAllStructs(b *testing.B) {
// 	structs := []struct {
// 		name  string
// 		frame Action
// 	}{
// 		// {"ReversedFrame", &ReversedFrame{}},
// 		// {"ShadowFrame", &ShadowFrame{}},
// 		// {"StealthFrame", &StealthFrame{}},
// 		// {"XDRFrame", &XDRFrame{}},
// 		// {"ShuffledFrame", &ShuffledFrame{}},
// 		// {"XORMaskedFrame", &XORMaskedFrame{}},
// 		// {"BinaryFrame", &BinaryFrame{}},
// 		{"NormalFrame", &NormalFrame{}},
// 	}
// 	testMsg := "Hello, World!"
// 	testTransport := byte(0x01)

// 	for _, s := range structs {
// 		b.Run(fmt.Sprintf("Benchmarking %s", s.name), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				encoded, err := s.frame.Encode(testMsg, testTransport)
// 				if err != nil {
// 					b.Fatalf("Encode failed for %s: %v", s.name, err)
// 				}
// 				_, _, _, err = s.frame.Decode(encoded)
// 				if err != nil {
// 					b.Fatalf("Decode failed for %s: %v", s.name, err)
// 				}
// 			}
// 		})
// 	}
// }
