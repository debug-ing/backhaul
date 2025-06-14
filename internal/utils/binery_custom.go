package utils

import (
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"
)

type Action interface {
	Encode(msg string, transport byte) ([]byte, error)
	EncodeString(msg string) ([]byte, error)

	Decode(data []byte) (string, byte, error)
	DecodeString(data []byte) (string, []byte, error)
}

func NewAction(name string) Action {
	switch name {
	case "NormalFrame":
		return &NormalFrame{}
	case "BinaryFrame":
		return &BinaryFrame{}
	// case "ReversedFrame":
	// 	return &ReversedFrame{}, nil
	// case "ShadowFrame":
	// 	return &ShadowFrame{}, nil
	// case "StealthFrame":
	// 	return &StealthFrame{}, nil
	// case "XDRFrame":
	// 	return &XDRFrame{}, nil
	// case "ShuffledFrame":
	// 	return &ShuffledFrame{}, nil
	// case "XORMaskedFrame":
	// 	return &XORMaskedFrame{}, nil
	default:
		return nil
	}
}

// -------------------- Normal Frame --------------------
type NormalFrame struct{}

func (n *NormalFrame) Encode(msg string, transport byte) ([]byte, error) {
	const headerSize = 3
	buf := make([]byte, headerSize+len(msg))
	binary.BigEndian.PutUint16(buf[:headerSize], uint16(len(msg)))
	buf[2] = transport
	copy(buf[headerSize:], msg)
	return buf, nil
}
func (n *NormalFrame) Decode(data []byte) (string, byte, error) {
	// // if len(data) < 3 {
	// // 	return "", nil, 0, fmt.Errorf("frame too short")
	// // }
	// //messageLength := binary.BigEndian.Uint16(lenBuf[:2])
	// messageLength := binary.BigEndian.Uint16(data[:])

	// payload := make([]byte, messageLength)
	// transport := payload[0]
	// return string(payload[1:]), payload, transport, nil
	if len(data) < 3 { // Minimum size check
		return "", 0, fmt.Errorf("frame too short")
	}

	// Extract message length
	messageLength := binary.BigEndian.Uint16(data[:2])
	if len(data) < int(3+messageLength) { // Check if data matches the expected length
		return "", 0, fmt.Errorf("payload length mismatch")
	}

	// Extract transport and payload
	transport := data[2]
	payload := data[3 : 3+messageLength]

	return string(payload), transport, nil
}

func (n *NormalFrame) EncodeString(msg string) ([]byte, error) {
	const headerSize = 2
	buf := make([]byte, headerSize+len(msg))
	binary.BigEndian.PutUint16(buf[:headerSize], uint16(len(msg)))
	copy(buf[headerSize:], msg)
	return buf, nil
}
func (n *NormalFrame) DecodeString(data []byte) (string, []byte, error) {
	if len(data) < 2 {
		return "", nil, fmt.Errorf("data too short to decode")
	}
	messageLength := binary.BigEndian.Uint16(data[:2])
	if len(data) < int(2+messageLength) {
		return "", nil, fmt.Errorf("data length mismatch")
	}
	messageBuf := data[2 : 2+messageLength]
	return string(messageBuf), messageBuf, nil
}

// -------------------- Normal Frame --------------------

// -------------------- Binary Frame --------------------
var maskTable [256]byte

type BinaryFrame struct {
}
type BinaryFrameData struct {
	Seed    byte
	Header  byte // encoded transport
	Length  byte // encoded length
	Payload []byte
}

func (b *BinaryFrame) Encode(msg string, transport byte) ([]byte, error) {
	// msgBytes := *(*[]byte)(unsafe.Pointer(&msg)) // no alloc conversion
	// l := len(msgBytes)

	// seed := byte(l*29+int(transport)) ^ 0xAA
	// rot := ((seed << 4) | (seed >> 4))

	// // Preallocate exact buffer
	// buf := make([]byte, 3+l)
	// buf[0] = seed
	// buf[1] = transport ^ seed
	// buf[2] = byte(l) ^ rot

	// maskIndex := int(seed)
	// for i := 0; i < l; i++ {
	// 	buf[3+i] = msgBytes[i] ^ maskTable[maskIndex]
	// 	maskIndex++
	// }
	// return buf, nil
	// const maxMessageLength = 1024 // حداکثر طول پیام
	// msgBytes := []byte(msg)
	// l := len(msgBytes)

	// // بررسی طول پیام
	// if l > maxMessageLength {
	// 	return nil, fmt.Errorf("message too long")
	// }

	// // تولید seed
	// seed := byte(42) // مقدار ثابت برای مثال؛ می‌توانید از مقدار تصادفی استفاده کنید
	// rot := ((seed << 4) | (seed >> 4))

	// // استفاده از آرایه ثابت برای جلوگیری از تخصیص حافظه
	// var buf [3 + maxMessageLength]byte
	// buf[0] = seed
	// buf[1] = transport ^ seed
	// buf[2] = byte(l) ^ rot

	// // رمزگذاری پیام
	// maskIndex := int(seed)
	// for i := 0; i < l; i++ {
	// 	buf[3+i] = msgBytes[i] ^ maskTable[maskIndex]
	// 	maskIndex++
	// }

	// // بازگرداندن داده‌های رمزگذاری‌شده
	// return buf[:3+l], nil

	const headerSize = 3
	msgBytes := []byte(msg) // تبدیل پیام به بایت‌ها
	l := len(msgBytes)

	// بررسی طول پیام
	if l > 65535 { // محدودیت طول پیام
		return nil, fmt.Errorf("message too long")
	}

	// تولید seed
	seed := byte(l*29+int(transport)) ^ 0xAA
	rot := (seed << 4) | (seed >> 4)

	// تخصیص بافر برای هدر و پیام
	buf := make([]byte, headerSize+l)
	buf[0] = seed
	buf[1] = transport ^ seed
	buf[2] = byte(l) ^ rot

	// رمزگذاری پیام
	maskIndex := int(seed)
	for i := 0; i < l; i++ {
		buf[headerSize+i] = msgBytes[i] ^ maskTable[maskIndex%len(maskTable)]
		maskIndex++
	}

	return buf, nil
}

func (b *BinaryFrame) EncodeString(msg string) ([]byte, error) {
	msgBytes := *(*[]byte)(unsafe.Pointer(&msg)) // no alloc conversion
	l := len(msgBytes)

	seed := byte(l*29+int(0)) ^ 0xAA // transport is 0 for string encoding
	rot := ((seed << 4) | (seed >> 4))

	// Preallocate exact buffer
	buf := make([]byte, 3+l)
	buf[0] = seed
	buf[1] = 0 ^ seed // transport is 0
	buf[2] = byte(l) ^ rot

	maskIndex := int(seed)
	for i := 0; i < l; i++ {
		buf[3+i] = msgBytes[i] ^ maskTable[maskIndex]
		maskIndex++
	}
	return buf, nil
}

func (b *BinaryFrame) Decode(buf []byte) (string, byte, error) {
	// if len(buf) < 3 {
	// 	return "", 0, fmt.Errorf("frame too short")
	// }

	// seed := buf[0]
	// transport := buf[1] ^ seed
	// rot := ((seed << 4) | (seed >> 4))
	// l := int(buf[2] ^ rot)

	// if len(buf) < 3+l {
	// 	return "", 0, fmt.Errorf("length mismatch")
	// }

	// out := make([]byte, l)
	// maskIndex := int(seed)
	// for i := 0; i < l; i++ {
	// 	out[i] = buf[3+i] ^ maskTable[maskIndex]
	// 	maskIndex++
	// }

	// return string(out[:l]), transport, nil
	// const maxMessageLength = 1024 // حداکثر طول پیام

	// // بررسی طول داده‌ها
	// if len(buf) < 3 {
	// 	return "", 0, fmt.Errorf("frame too short")
	// }

	// // استخراج هدر
	// seed := buf[0]
	// transport := buf[1] ^ seed
	// rot := ((seed << 4) | (seed >> 4))
	// l := int(buf[2] ^ rot)

	// // بررسی طول پیام
	// if l > maxMessageLength {
	// 	return "", 0, fmt.Errorf("message too long")
	// }
	// if len(buf) < 3+l {
	// 	return "", 0, fmt.Errorf("length mismatch")
	// }

	// // استفاده از آرایه ثابت برای جلوگیری از تخصیص حافظه
	// var out [maxMessageLength]byte
	// maskIndex := int(seed)
	// for i := 0; i < l; i++ {
	// 	out[i] = buf[3+i] ^ maskTable[maskIndex]
	// 	maskIndex++
	// }

	// // بازگرداندن پیام رمزگشایی‌شده
	// return string(out[:l]), transport, nil

	// const headerSize = 3

	// // بررسی حداقل طول داده‌ها
	// if len(buf) < headerSize {
	// 	return "", 0, fmt.Errorf("frame too short")
	// }

	// // استخراج طول پیام از هدر
	// messageLength := binary.BigEndian.Uint16(buf[:2])

	// // بررسی طول داده‌ها
	// if len(buf) < int(headerSize+messageLength) {
	// 	return "", 0, fmt.Errorf("payload length mismatch")
	// }

	// // استخراج transport و پیام
	// transport := buf[2]
	// message := string(buf[headerSize : headerSize+messageLength])

	// return message, transport, nil

	const headerSize = 3

	// بررسی حداقل طول داده‌ها
	if len(buf) < headerSize {
		return "", 0, fmt.Errorf("frame too short")
	}

	// استخراج هدر
	seed := buf[0]
	transport := buf[1] ^ seed
	rot := (seed << 4) | (seed >> 4)
	messageLength := int(buf[2] ^ rot)

	// بررسی طول پیام
	if len(buf) < headerSize+messageLength {
		return "", 0, fmt.Errorf("payload length mismatch")
	}

	// رمزگشایی پیام
	payload := buf[headerSize : headerSize+messageLength]
	decodedMessage := make([]byte, messageLength)
	maskIndex := int(seed)
	for i := 0; i < messageLength; i++ {
		decodedMessage[i] = payload[i] ^ maskTable[maskIndex%len(maskTable)]
		maskIndex++
	}

	// بازگرداندن پیام رمزگشایی‌شده
	return string(decodedMessage), transport, nil
}

func (b *BinaryFrame) DecodeString(buf []byte) (string, []byte, error) {
	if len(buf) < 3 {
		return "", nil, fmt.Errorf("frame too short")
	}

	seed := buf[0]
	rot := ((seed << 4) | (seed >> 4))
	l := int(buf[2] ^ rot)

	if len(buf) < 3+l {
		return "", nil, fmt.Errorf("length mismatch")
	}

	out := make([]byte, l)
	maskIndex := int(seed)
	for i := 0; i < l; i++ {
		out[i] = buf[3+i] ^ maskTable[maskIndex]
		maskIndex++
	}

	return bytesToString(out), out, nil
}
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// -------------------- Binary Frame --------------------

// -------------------- Reversed Frame --------------------
type ReversedFrame struct{}

func (r *ReversedFrame) Encode(msg string, transport byte) ([]byte, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}
	start := byte(0xAB)
	end := byte(0xCD)
	reversed := []byte(msg)
	for i := 0; i < len(reversed)/2; i++ {
		reversed[i], reversed[len(reversed)-1-i] = reversed[len(reversed)-1-i], reversed[i]
	}
	frame := []byte{start}
	frame = append(frame, reversed...)
	frame = append(frame, transport, end)
	return frame, nil
}
func (r *ReversedFrame) Decode(data []byte) (string, byte, error) {
	if len(data) < 3 || data[0] != 0xAB || data[len(data)-1] != 0xCD {
		return "", 0, fmt.Errorf("invalid frame")
	}
	transport := data[len(data)-2]
	reversed := data[1 : len(data)-2]
	for i := 0; i < len(reversed)/2; i++ {
		reversed[i], reversed[len(reversed)-1-i] = reversed[len(reversed)-1-i], reversed[i]
	}
	return string(reversed), transport, nil
}

// -------------------- Reversed Frame --------------------

// -------------------- Shadow Frame --------------------
type ShadowFrame struct{}

func (s *ShadowFrame) Encode(msg string, transport byte) ([]byte, error) {
	const maxLength = 31
	if len(msg) > maxLength {
		return nil, fmt.Errorf("message too long (max %d bytes)", maxLength)
	}
	seed := byte(time.Now().UnixNano() % 256)
	key := seed * 73
	meta := (transport & 0x07) << 5
	meta |= byte(len(msg)) & 0x1F
	encrypted := make([]byte, len(msg))
	for i := range msg {
		encrypted[i] = msg[i] ^ key
	}
	frame := append([]byte{seed, meta}, encrypted...)
	return frame, nil
}

func (s *ShadowFrame) Decode(data []byte) (string, byte, error) {
	if len(data) < 2 {
		return "", 0, fmt.Errorf("invalid frame")
	}
	seed := data[0]
	meta := data[1]
	key := seed * 73

	transport := (meta >> 5) & 0x07
	length := int(meta & 0x1F)

	if len(data) < 2+length {
		return "", 0, fmt.Errorf("frame too short")
	}

	payload := data[2 : 2+length]
	decrypted := make([]byte, length)
	for i := range payload {
		decrypted[i] = payload[i] ^ key
	}

	return string(decrypted), transport, nil
}

// -------------------- Shadow Frame --------------------

// -------------------- Stealth Frame --------------------

type StealthFrame struct{}

func (s *StealthFrame) Encode(msg string, transport byte) ([]byte, error) {
	const maxLength = 31
	if len(msg) > maxLength {
		return nil, fmt.Errorf("message too long (max %d bytes)", maxLength)
	}
	seed := byte(time.Now().UnixNano() % 256)
	key := seed * 73
	meta := (transport & 0x07) << 5
	meta |= byte(len(msg)) & 0x1F
	encrypted := make([]byte, len(msg))
	for i := range msg {
		encrypted[i] = msg[i] ^ key
	}
	frame := append([]byte{seed, meta}, encrypted...)
	return frame, nil
}

func (s *StealthFrame) Decode(data []byte) (string, byte, error) {
	if len(data) < 2 {
		return "", 0, fmt.Errorf("invalid frame")
	}
	seed := data[0]
	meta := data[1]
	key := seed * 73

	transport := (meta >> 5) & 0x07
	length := int(meta & 0x1F)

	if len(data) < 2+length {
		return "", 0, fmt.Errorf("frame too short")
	}

	payload := data[2 : 2+length]
	decrypted := make([]byte, length)
	for i := range payload {
		decrypted[i] = payload[i] ^ key
	}

	return string(decrypted), transport, nil
}

// TODO check this
// func EncodeMinimal(msg string, transport byte) ([]byte, error) {
// 	seed := byte(len(msg)*31+int(transport)) & 0xFF
// 	rot := ((seed << 3) | (seed >> 5)) & 0xFF

// 	length := byte(len(msg)) ^ rot
// 	header := transport ^ seed

// 	payload := make([]byte, len(msg))
// 	for i := range msg {
// 		// mask := (seed * byte(i+1)) % 256
// 		mask := byte(int(seed) * (i + 1) % 256)
// 		payload[i] = msg[i] ^ mask
// 	}

// 	frame := append([]byte{seed, header, length}, payload...)
// 	return frame, nil
// }

// func DecodeMinimal(frame []byte) (msg string, transport byte, err error) {
// 	if len(frame) < 3 {
// 		err = fmt.Errorf("frame too short")
// 		return
// 	}

// 	seed := frame[0]
// 	rot := ((seed << 3) | (seed >> 5)) & 0xFF
// 	transport = frame[1] ^ seed
// 	length := frame[2] ^ rot

// 	if len(frame[3:]) != int(length) {
// 		err = fmt.Errorf("payload length mismatch")
// 		return
// 	}

// 	payload := make([]byte, length)
// 	for i := range payload {
// 		// mask := (seed * byte(i+1)) % 256
// 		mask := byte(int(seed) * (i + 1) % 256)
// 		payload[i] = frame[3+i] ^ mask
// 	}

// 	msg = string(payload)
// 	return
// }

// -------------------- xdr Frame --------------------
type XDRFrame struct{}

func (x *XDRFrame) Encode(msg string, transport byte) ([]byte, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}
	mask := byte(0x7A ^ byte(len(msg)))
	length := byte(len(msg)) ^ mask
	transportMasked := transport ^ mask

	data := make([]byte, len(msg))
	for i, c := range []byte(msg) {
		data[i] = c ^ mask
	}

	frame := []byte{mask, length, transportMasked}
	return append(frame, data...), nil
}
func (x *XDRFrame) Decode(data []byte) (string, byte, error) {
	if len(data) < 3 {
		return "", 0, fmt.Errorf("frame too short")
	}
	mask := data[0]
	length := data[1] ^ mask
	transport := data[2] ^ mask

	if len(data) < int(length)+3 {
		return "", 0, fmt.Errorf("payload length mismatch")
	}

	payload := make([]byte, length)
	for i := range payload {
		payload[i] = data[3+i] ^ mask
	}

	return string(payload), transport, nil
}

// -------------------- xdr Frame --------------------

// -------------------- shuffled Frame --------------------
type ShuffledFrame struct{}

func (s *ShuffledFrame) Encode(msg string, transport byte) ([]byte, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}
	pos := byte(len(msg) / 2)
	payload := []byte(msg)
	payload = append(payload[:pos], append([]byte{transport}, payload[pos:]...)...)
	return append([]byte{byte(len(payload)), pos}, payload...), nil
}

func (s *ShuffledFrame) Decode(data []byte) (string, byte, error) {
	length := int(data[0])
	pos := int(data[1])
	payload := data[2 : 2+length]
	transport := payload[pos]
	msg := append(payload[:pos], payload[pos+1:]...)
	return string(msg), transport, nil
}

// -------------------- shuffled Frame --------------------

// -------------------- XOR Masked Frame --------------------
type XORMaskedFrame struct{}

func (x *XORMaskedFrame) Encode(msg string, transport byte) ([]byte, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}
	mask := byte(0x5A ^ byte(len(msg)))
	length := byte(len(msg)) ^ mask
	transportMasked := transport ^ mask

	data := make([]byte, len(msg))
	for i, c := range []byte(msg) {
		data[i] = c ^ mask
	}

	frame := []byte{mask, length, transportMasked}
	return append(frame, data...), nil
}
func (x *XORMaskedFrame) Decode(data []byte) (string, byte, error) {
	if len(data) < 3 {
		return "", 0, fmt.Errorf("frame too short")
	}
	mask := data[0]
	length := data[1] ^ mask
	transport := data[2] ^ mask

	if len(data) < int(length)+3 {
		return "", 0, fmt.Errorf("payload length mismatch")
	}

	payload := make([]byte, length)
	for i := range payload {
		payload[i] = data[3+i] ^ mask
	}

	return string(payload), transport, nil
}

// -------------------- XOR Masked Frame --------------------
