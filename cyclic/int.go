////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package cyclic wraps our large.Int structure.  It is designed to be used in
// conjunction with the cyclic.Group object. The cyclic.Group object
// will provide implementations of various modular operations within the group.
// A cyclic.IntBuffer type will be created to store large batches of groups.
package cyclic

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/large"
)

// Create the cyclic.Int type as a wrapper of a large.Int
// and group fingerprint
type Int struct {
	value       *large.Int
	fingerprint uint64
}

// ByteLen gets the byte length of cyclic int
func (z *Int) ByteLen() int {
	byteLen := z.value.ByteLen()
	return byteLen
}

// GetLargeInt gets a deepcopy of the largeInt from cyclicInt
// This is necessary because otherwise the internal
// value of the into could be edited and made to be
// outside the group.
func (z *Int) GetLargeInt() *large.Int {
	r := large.NewInt(0)
	r.Set(z.value)
	return r
}

// GetGroupFingerprint gets the group fingerprint from cyclicInt
func (z *Int) GetGroupFingerprint() uint64 {
	return z.fingerprint
}

// Bits gets the underlying word slice of cyclic int
// Use this for low-level functions where speed is critical
// For speed reasons, I don't copy here. This could allow the int to be set outside of the group
func (z *Int) Bits() large.Bits {
	return z.value.Bits()
}

// Bytes gets the bytes of cyclicInt value
func (z *Int) Bytes() []byte {
	return z.value.Bytes()
}

// LeftpadBytes gets left padded bytes of cyclicInt value
func (z *Int) LeftpadBytes(length uint64) []byte {
	return z.value.LeftpadBytes(length)
}

// BitLen gets the length of the cyclic int
func (z *Int) BitLen() int {
	return z.value.BitLen()
}

// DeepCopy returns a complete copy of the cyclic int such that no
// underlying data is linked
func (z *Int) DeepCopy() *Int {
	return &Int{
		z.value.DeepCopy(),
		z.fingerprint,
	}
}

// Compare two cyclicInts
// returns -2 if fingerprint differs
// returns value.Cmp otherwise
func (z *Int) Cmp(x *Int) int {
	if z.fingerprint != x.fingerprint {
		return -2
	}
	return z.value.Cmp(x.value)
}

// Reset cyclicInt to 1
func (z *Int) Reset() {
	z.value.SetInt64(1)
}

// Return truncated base64 encoded string of group fingerprint
func (z *Int) textFingerprint(length int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, z.fingerprint)
	fullText := base64.StdEncoding.EncodeToString(buf)
	if length == 0 || len(fullText) <= length {
		return fullText
	} else {
		return fullText[:length] + "..."
	}
}

// Text returns the string representation of z in the given base. Base
// must be between 2 and 36, inclusive. The result uses the lower-case
// letters 'a' to 'z' for digit values >= 10. No base prefix (such as
// "0x") is added to the string.
// Text truncates ints to a length of 10, appending an ellipsis
// if the int is too long.
// The group fingerprint is base64 encoded and also truncated
// z is then represented as: value... in GRP: fingerprint...
func (z *Int) Text(base int) string {
	const intTextLen = 10
	return z.TextVerbose(base, intTextLen)
}

// TextVerbose returns the string representation of z in the given base. Base
// must be between 2 and 36, inclusive. The result uses the lower-case
// letters 'a' to 'z' for digit values >= 10. No base prefix (such as
// "0x") is added to the string.
// TextVerbose truncates ints to a length of length in characters (not runes)
// and append an ellipsis to indicate that the whole int wasn't returned,
// unless len is 0, in which case it will return the whole int as a string.
// The group fingerprint is base64 encoded and also truncated
// z is then represented as: value... in GRP: fingerprint...
func (z *Int) TextVerbose(base int, length int) string {
	valueText := z.value.TextVerbose(base, length)
	fingerprintText := z.textFingerprint(length)
	return valueText + " in GRP: " + fingerprintText
}

// GOB decode bytes to cyclicInt
func (z *Int) GobDecode(in []byte) error {
	// anonymous structure
	s := struct {
		F []byte
		V []byte
	}{
		make([]byte, 8),
		[]byte{},
	}

	var buf bytes.Buffer

	// Write bytes to the buffer
	buf.Write(in)

	// Create new decoder that reads from the buffer
	dec := gob.NewDecoder(&buf)

	// Receive and decode data
	err := dec.Decode(&s)

	if err != nil {
		return err
	}

	// Convert decoded bytes and put into empty structure
	z.value = large.NewIntFromBytes(s.V)
	z.fingerprint = binary.BigEndian.Uint64(s.F)

	return nil
}

// GOB encode cyclicInt to bytes
func (z *Int) GobEncode() ([]byte, error) {
	// Anonymous structure
	s := struct {
		F []byte
		V []byte
	}{
		make([]byte, 8),
		z.Bytes(),
	}

	binary.BigEndian.PutUint64(s.F, z.fingerprint)
	var buf bytes.Buffer

	// Create new encoder that will transmit the buffer
	enc := gob.NewEncoder(&buf)

	// Transmit the data
	err := enc.Encode(s)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// BinaryEncode encodes the Int into a compressed byte format.
func (z *Int) BinaryEncode() []byte {
	var buff bytes.Buffer
	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, z.fingerprint)
	buff.Write(b)
	buff.Write(z.Bytes())

	return buff.Bytes()
}

// BinaryDecode decompresses the encoded byte slice to an Int.
func (z *Int) BinaryDecode(b []byte) error {
	if len(b) < 8 {
		return errors.Errorf("length of buffer %d != %d expected", len(b), 8)
	}

	buff := bytes.NewBuffer(b)
	z.fingerprint = binary.LittleEndian.Uint64(buff.Next(8))
	z.value = large.NewIntFromBytes(buff.Bytes())
	return nil
}

// Erase overwrite all underlying data from a cyclic Int by setting its value
// and fingerprint to zero. All underlying released data will be removed by the
// garbage collector.
func (z *Int) Erase() {
	z.value.SetInt64(0)
	z.fingerprint = 0
}

// -------------- Marshal Operators -------------- //
// intData holds the value of a cyclic int in public fields to allow for
// marshalling and unmarshalling.
type intData struct {
	Value       *large.Int
	Fingerprint uint64
}

// MarshalJSON is a custom marshaling function for cyclic int. It is used when
// json.Marshal is called on a large int.
func (z *Int) MarshalJSON() ([]byte, error) {
	data := intData{
		Value:       z.value,
		Fingerprint: z.fingerprint,
	}

	return json.Marshal(data)
}

// UnmarshalJSON is a custom unmarshalling function for cyclic int. It is used
// when json.Unmarshal is called on a large int.
func (z *Int) UnmarshalJSON(b []byte) error {
	data := &intData{}
	err := json.Unmarshal(b, data)
	if err != nil {
		return err
	}

	z.value = data.Value
	z.fingerprint = data.Fingerprint

	return nil
}
