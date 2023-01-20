////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package cyclic

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"gitlab.com/xx_network/crypto/large"
	"math/rand"
	"reflect"
	"testing"
)

// Tests NewGroup functionality
func TestNewGroup(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	actual := NewGroup(p, g)

	type testStruct struct {
		prime *large.Int
		g     *large.Int
	}
	expected := testStruct{p, g}
	if actual.prime.Cmp(expected.prime) != 0 {
		t.Errorf("TestNewGroup failed to initialize prime, expected: '%v',"+
			" got: '%v'", expected.prime.Text(10), actual.prime.Text(10))
	} else if actual.gen.Cmp(expected.g) != 0 {
		t.Errorf("TestNewGroup failed to initialize generator, expected: '%v',"+
			" got: '%v'", expected.g.Text(10), actual.gen.Text(10))
	}
}

// Test creation of cyclicInt in the group from int64
func TestNewInt(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := large.NewInt(42)
	actual := grp.NewInt(42)

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("NewInt creation failed, expected: %v,"+
			"got: %v", expected, actual.value)
	} else if actual.GetGroupFingerprint() != grp.GetFingerprint() {
		t.Errorf("NewInt is not in the group, expected group fingerprint: %v,"+
			"got: %v", grp.GetFingerprint(), actual.GetGroupFingerprint())
	}
}

//Tests creation and properties of an IntBuffer
func TestGroup_NewIntBuffer(t *testing.T) {

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	//test that the size is correct and the default value is set correctly
	rng := rand.New(rand.NewSource(42))

	tests := 100

	for i := 0; i < tests; i++ {
		defaultInt := grp.Random(grp.NewInt(1))
		size := rng.Uint32() % 10000
		buf := grp.NewIntBuffer(size, defaultInt)

		//test that the length is correct
		if len(buf.values) != int(size) {
			t.Errorf("NewIntBuffer did not generate buffer of the correct size: "+
				"Expected %v, Recieved: %v", size, len(buf.values))
		}

		pass := true

		defaultIntLarge := defaultInt.GetLargeInt()

		//test that the default value is set correctly
		for _, i := range buf.values {
			if i.Cmp(defaultIntLarge) != 0 {
				pass = false
			}
		}

		if !pass {
			t.Errorf("NewIntBuffer internal values not equal to default value")
		}
	}

	//test that when passed default int is nil values are set to prime-1
	buf := grp.NewIntBuffer(10, nil)

	for _, i := range buf.values {
		if i.Cmp(grp.psub1) != 0 {
			t.Errorf("NewIntBuffer internal values not equal to psub1 when nill passed")
		}
	}

}

// Test creation of cyclicInt in the group from int64 fails when outside the group
func TestNewInt_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	grp.NewInt(0)

	t.Errorf("NewInt created even when outside of the group")
}

// Test creation of cyclicInt in the group from large.Int
func TestNewIntFromLargeInt(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := large.NewInt(42)
	actual := grp.NewIntFromLargeInt(expected)

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("NewIntFromLargeInt creation failed, expected: %v,"+
			"got: %v", expected, actual.value)
	} else if actual.GetGroupFingerprint() != grp.GetFingerprint() {
		t.Errorf("NewIntFromLargeInt is not in the group, expected group fingerprint: %v,"+
			"got: %v", grp.GetFingerprint(), actual.GetGroupFingerprint())
	}
}

// Test creation of cyclicInt in the group from large.Int fails when outside the group
func TestNewIntFromLargeInt_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	grp.NewIntFromLargeInt(large.NewInt(0))

	t.Errorf("NewIntFromLargeInt created even when outside of the group")
}

// Test creation of cyclicInt in the group from byte array
func TestNewIntFromBytes(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := large.NewInt(42)
	value := []byte{0x2A}
	actual := grp.NewIntFromBytes(value)

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("NewIntFromBytes creation failed, expected: %v,"+
			"got: %v", expected, actual.value)
	} else if actual.GetGroupFingerprint() != grp.GetFingerprint() {
		t.Errorf("NewIntFromBytes is not in the group, expected group fingerprint: %v,"+
			"got: %v", grp.GetFingerprint(), actual.GetGroupFingerprint())
	}
}

// Test creation of cyclicInt in the group from bytes fails when outside the group
func TestNewIntFromBytes_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	grp.NewIntFromBytes([]byte{0})

	t.Errorf("NewIntFromBytes created even when outside of the group")
}

// Test creation of cyclicInt in the group from string
// Also confirm that if the string can't be converted, nil is returned
func TestNewIntFromString(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := large.NewInt(42)
	value := "42"
	actual := grp.NewIntFromString(value, 10)

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("NewIntFromString creation failed, expected: %v,"+
			"got: %v", expected, actual.value)
	} else if actual.GetGroupFingerprint() != grp.GetFingerprint() {
		t.Errorf("NewIntFromString is not in the group, expected group fingerprint: %v,"+
			"got: %v", grp.GetFingerprint(), actual.GetGroupFingerprint())
	}

	errVal := grp.NewIntFromString("185", 5)

	if errVal != nil {
		t.Errorf("NewIntFromString should return nil when error occurs decoding string")
	}
}

// Test creation of cyclicInt in the group from string fails when outside the group
func TestNewIntFromString_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	grp.NewIntFromString("0", 16)

	t.Errorf("NewIntFromString created even when outside of the group")
}

// Show that NewIntFromBits fails when outside of group
func TestGroup_NewIntFromBits_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	grp.NewIntFromBits(large.Bits{0})

	t.Errorf("NewIntFromBits created even when outside of the group")
}

// Show that NewIntFromBits makes a big int from a word string
func TestGroup_NewIntFromBits(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := grp.NewIntFromString("123456", 16)
	i := grp.NewIntFromBits(large.Bits{0x123456})
	t.Log(i.TextVerbose(16, 0))

	if expected.Cmp(i) != 0 {
		t.Errorf("Expected int to be %v, got %v", expected.Text(16), i.Text(16))
	}
}

// Test creation of cyclicInt in the group from Max4KInt value
func TestNewMaxInt(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := grp.psub1
	actual := grp.NewMaxInt()

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("NewMaxInt creation failed, expected: %v,"+
			"got: %v", expected, actual.value)
	} else if actual.GetGroupFingerprint() != grp.GetFingerprint() {
		t.Errorf("NewMaxInt is not in the group, expected group fingerprint: %v,"+
			"got: %v", grp.GetFingerprint(), actual.GetGroupFingerprint())
	}
}

// Test creation of cyclicInt in the group from uint64
func TestNewIntFromUInt(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := large.NewInt(42)
	actual := grp.NewIntFromUInt(uint64(42))

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("NewIntFromUInt creation failed, expected: %v,"+
			"got: %v", expected, actual.value)
	} else if actual.GetGroupFingerprint() != grp.GetFingerprint() {
		t.Errorf("NewIntFromUInt is not in the group, expected group fingerprint: %v,"+
			"got: %v", grp.GetFingerprint(), actual.GetGroupFingerprint())
	}
}

// Test creation of cyclicInt in the group from uint64 fails when outside the group
func TestNewIntFromUInt_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	grp.NewIntFromUInt(0)

	t.Errorf("NewIntFromUInt created even when outside of the group")
}

// Test group fingerprint getter
func TestGetFingerprint(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	h := sha256.New()
	h.Write(p.Bytes())
	h.Write(g.Bytes())
	expected := large.NewIntFromBytes(h.Sum(nil)[:GroupFingerprintSize]).Uint64()

	if grp.GetFingerprint() != expected {
		t.Errorf("GetFingerprint returned wrong value, expected: %v,"+
			" got: %v", expected, grp.GetFingerprint())
	}
}

// Test group fingerprint getter
func TestGetFingerprintText(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := "ln9lzlk2..."

	if grp.GetFingerprintText() != expected {
		t.Errorf("GetFingerprintText returned wrong value, expected: %v,"+
			" got: %v", expected, grp.GetFingerprintText())
	}
}

// Test setting cyclicInt to another from the same group
func TestSet(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := grp.NewInt(int64(42))
	actual := grp.NewInt(int64(69))

	if actual.Cmp(expected) == 0 {
		t.Errorf("Original Ints should be different to test Set")
	}

	grp.Set(actual, expected)

	result := actual.Cmp(expected)

	if result != 0 {
		t.Errorf("Test of Set failed, expected: '0', got: '%v'",
			result)
	}
}

// Test that Set panics if cyclicInt doesn't belong to the group
func TestSet_Panic(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	expected := grp.NewInt(int64(42))
	actual := grp2.NewInt(int64(69))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Set should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	grp.Set(actual, expected)
}

// TestSetLargeInt tests if the value is properly copied
func TestSetLargeInt(t *testing.T) {
	p := large.NewInt(17)
	g := large.NewInt(7)
	group := NewGroup(p, g)

	inputs := []int64{2, 1, 12, 17, 18}

	for i := 0; i < len(inputs); i++ {
		tmp := group.NewInt(5)
		v := tmp.value
		exp := group.NewInt(inputs[i%3])
		li := large.NewInt(inputs[i])
		if group.SetLargeInt(tmp, li) != nil {
			if i >= 3 {
				t.Errorf("TestLargeInt set value outside "+
					"of group: %d", inputs[i])
			}
		}
		if i < 3 && (tmp.Cmp(exp) != 0 || v.Cmp(li) != 0) {
			t.Errorf("TestSetFromLargeInt failed at index %v", i)
		}
	}
}

// Test setting cyclicInt in the same group from bytes
func TestSetBytes(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := []*Int{
		grp.NewInt(42),
		grp.NewInt(6553522),
		grp.NewInt(2)}
	testBytes := [][]byte{
		{0x2A},             // 42
		{0x63, 0xFF, 0xB2}, // 6553522
		{0x02}}

	actual := grp.NewInt(55)

	for i, testi := range testBytes {
		actual = grp.SetBytes(actual, testi)
		if actual.Cmp(expected[i]) != 0 {
			t.Errorf("Test of SetBytes failed at index %v, expected: '%v', "+
				"actual: %v", i, expected[i].Text(10), actual.Text(10))
		}
	}

}

// Test that SetBytes panics if cyclicInt doesn't belong to the group
func TestSetBytes_Panic(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	actual := grp2.NewInt(int64(42))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetBytes should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	grp.SetBytes(actual, []byte("TEST"))
}

// Shows that setting bits results in a different integer
func TestGroup_SetBits(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := grp.NewInt(2)
	newInt := grp.NewInt(1)
	grp.SetBits(newInt, large.Bits{2})
	if expected.Cmp(newInt) != 0 {
		t.Errorf("Setbits didn't set to the expected int result. Got %v, expected %v", newInt.Text(16), expected.Text(16))
	}
}

// Shows that setting bits with a different group from the original integer panics
func TestGroup_SetBits_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	i := grp.NewInt(1)
	grp2.SetBits(i, large.Bits{123456})

	t.Errorf("SetBits worked even when another group was used")
}

// OverwriteBits is a higher-level method that copies a bits slice into an integer
// or allocates a new slice if necessary
// Shows that OverwriteBits never uses the passed bits slice to back the new integer
// Shows that the results are as expected, even in cases where the lengths aren't equal
func TestGroup_OverwriteBits(t *testing.T) {
	// We need a bigger group than the other tests, because for this test to be meaningful
	// the prime needs to take up more than one word
	// Using modp1536
	p := large.NewIntFromString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA237327FFFFFFFFFFFFFFFF", 16)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	// Case 1: enough space in the existing int for the overwriting int to exist
	existing := grp.NewMaxInt()
	existingStart := &existing.Bits()[0]
	t.Log("existing cap:", cap(existing.Bits()))
	expected := grp.NewIntFromString("1234567890abcdef1234567890abcdef1234567890abcdef123", 16)
	t.Log("expected len:", len(expected.Bits()))
	grp.OverwriteBits(existing, expected.Bits())

	// Start of backing bits slice should be the same
	if existingStart != &existing.Bits()[0] {
		t.Errorf("start of existing changed. had %v, got %v", existingStart, &existing.Bits()[0])
	}
	// Should get expected number in existing
	if existing.Cmp(expected) != 0 {
		t.Errorf("actual differed from expected. actual: %v, expected %v", existing, expected)
	}

	// Case 2: not enough space in the existing int for the overwriting int to exist
	newAlloc := grp.NewInt(1)
	newAllocStart := &newAlloc.Bits()[0]
	t.Log("newAlloc cap:", cap(newAlloc.Bits()))
	grp.OverwriteBits(newAlloc, expected.Bits())

	// Start of backing bits slice should be different
	if newAllocStart == &newAlloc.Bits()[0] {
		t.Errorf("start of newAlloc didn't change. had %v, got %v", newAllocStart, &newAlloc.Bits()[0])
	}
	// It should also be different from the backing bits slice of expected
	if &newAlloc.Bits()[0] == &existing.Bits()[0] {
		t.Errorf("newAlloc uses same backing memory as existing! bad!!")
	}
	// Of course, the int should be equal to existing
	if expected.Cmp(newAlloc) != 0 {
		t.Errorf("actual differed from expected. actual: %v, expected %v", newAlloc, expected)
	}
}

// Is it possible to trigger the panic?
func TestGroup_OverwriteBits_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	i := grp.NewInt(1)
	grp2.OverwriteBits(i, large.Bits{123456})

	t.Errorf("OverwriteBits worked with a different group")
}

// Test setting cyclicInt in the same group from string
func TestSetString(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	type testStructure struct {
		str  string
		base int
	}

	testStructs := []testStructure{
		{"42", 0},
		{"100000000", 0},
		{"5", 0},
		{"1", 0},
		{"f", 0},
		{"182", 5},
		{"10", 2},
	}

	expected := []*Int{
		grp.NewInt(42),
		grp.NewInt(100000000),
		grp.NewInt(5),
		grp.NewInt(1),
		nil,
		nil,
		grp.NewInt(2),
	}

	actual := grp.NewInt(1)

	for i, testi := range testStructs {
		ret := grp.SetString(actual, testi.str, testi.base)

		// Test invalid input making sure it occurs when expected
		if ret == nil {
			if expected[i] != nil {
				t.Error("Test of SetString() failed at index:", i,
					"Function didn't handle invalid input correctly")
			}
		} else {
			if actual.Cmp(expected[i]) != 0 {
				t.Errorf("Test of SetString() failed at index: %v Expected: %v;"+
					" Actual: %v", i, expected[i].Text(10), actual.Text(10))
			}
		}
	}

}

// Test that SetString panics if cyclicInt doesn't belong to the group
func TestSetString_Panic(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	actual := grp2.NewInt(int64(42))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetString should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	grp.SetString(actual, "1234", 10)
}

// Test setting cyclicInt in the same group to Max4KInt value
func TestSetMaxInt(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := grp.GetPSub1()
	actual := grp.NewInt(int64(69))

	if actual.Cmp(expected) == 0 {
		t.Errorf("Original Ints should be different to test SetMaxInt")
	}

	grp.SetMaxInt(actual)

	result := actual.Cmp(expected)

	if result != 0 {
		t.Errorf("Test of SetMaxInt failed, expected: '0', got: '%v'",
			result)
	}
}

// Test that Set panics if cyclicInt doesn't belong to the group
func TestSetMaxInt_Panic(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	actual := grp2.NewInt(int64(69))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetMaxInt should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	grp.SetMaxInt(actual)
}

// Test setting cyclicInt in the same group to uint64 value
func TestSetUint64(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	expected := grp.NewInt(int64(42))
	actual := grp.NewInt(int64(69))

	if actual.Cmp(expected) == 0 {
		t.Errorf("Original Ints should be different to test SetUint64")
	}

	grp.SetUint64(actual, uint64(42))

	result := actual.Cmp(expected)

	if result != 0 {
		t.Errorf("Test of SetUint64 failed, expected: '0', got: '%v'",
			result)
	}
}

// Test that Set panics if cyclicInt doesn't belong to the group
func TestSetUint64_Panic(t *testing.T) {
	p := large.NewInt(1000000010101111111)
	g := large.NewInt(5)
	grp := NewGroup(p, g)
	g2 := large.NewInt(2)
	grp2 := NewGroup(p, g2)

	actual := grp2.NewInt(int64(69))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetUint64 should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	grp.SetUint64(actual, uint64(0))
}

// Test multiplication under the group
func TestMul(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)

	actual := []*Int{
		group.Mul(group.NewInt(20), group.NewInt(11), group.NewInt(1)),
		group.Mul(group.NewInt(1), group.NewInt(10), group.NewInt(1)),
	}
	expected := []*Int{
		group.NewInt((20 * 11) % prime),
		group.NewInt(10),
	}

	for i := 0; i < len(actual); i++ {
		if actual[i].value.Cmp(expected[i].value) != 0 {
			t.Errorf("TestMulForGroup failed at index:%v, expected:%v, got:%v",
				i, expected[i].value.Text(10), actual[i].value.Text(10))
		}
	}

}

// Test left padded bytes getter
func TestFullBytes(t *testing.T) {

	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2A}

	actual := grp.NewInt(int64(42))

	actualbytes := grp.FullBytes(actual)

	if !bytes.Equal(actualbytes, expected) {
		t.Errorf("Test of FullBytes failed, expected: '%v', got: '%v'",
			actualbytes, expected)
	}
}

// Test that mul panics if cyclicInt doesn't belong to the group
func TestMul_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group.NewInt(20)
	b := group2.NewInt(11)
	c := group.NewInt(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Mul should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	group.Mul(a, b, c)
}

// Test Inside that checks if a number is inside the group
func TestInside(t *testing.T) {
	p := large.NewInt(17)
	g := large.NewInt(7)
	group := NewGroup(p, g)
	expected := []bool{
		false,
		true,
		false,
		false,
		true,
	}
	actual := []bool{
		group.Inside(large.NewInt(0)),
		group.Inside(large.NewInt(1)),
		group.Inside(large.NewInt(17)),
		group.Inside(large.NewInt(18)),
		group.Inside(large.NewInt(12)),
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("TestInside failed at index:%v, expected:%v, got:%v",
				i, expected[i], actual[i])
		}
	}
}

// Test Inside that checks if a number is inside the group
func TestSingleBytesInside(t *testing.T) {
	p := large.NewInt(1023)
	g := large.NewInt(7)
	group := NewGroup(p, g)
	expected := []bool{
		false,
		false,
		true,
		true,
		true,
		true,
		false,
		false,
	}
	actual := []bool{
		group.singleBytesInside([]byte{0}),
		group.singleBytesInside(large.NewInt(0).Bytes()),
		group.singleBytesInside(large.NewInt(1).Bytes()),
		group.singleBytesInside(large.NewInt(17).Bytes()),
		group.singleBytesInside(large.NewInt(70).Bytes()),
		group.singleBytesInside(large.NewInt(1022).Bytes()),
		group.singleBytesInside(large.NewInt(1111).Bytes()),
		group.singleBytesInside(large.NewInt(100000).Bytes()),
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("TestBytesInside failed at index:%v, expected:%v, got:%v",
				i, expected[i], actual[i])
		}
	}
}

// Test Inside that checks if a number is inside the group
func TestBytesInside(t *testing.T) {
	p := large.NewInt(1023)
	g := large.NewInt(7)
	group := NewGroup(p, g)
	expected := []bool{
		true,
		true,
		false,
		false,
	}
	actual := []bool{
		group.BytesInside(large.NewInt(1).Bytes()),
		group.BytesInside(large.NewInt(1).Bytes(), large.NewInt(1000).Bytes(), large.NewInt(300).Bytes()),
		group.BytesInside(large.NewInt(1).Bytes(), large.NewInt(1000).Bytes(), large.NewInt(300).Bytes(), large.NewInt(2000).Bytes()),
		group.BytesInside(large.NewInt(0).Bytes(), large.NewInt(1100).Bytes(), large.NewInt(30000000).Bytes(), large.NewInt(400900).Bytes()),
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("TestBytesInside failed at index:%v, expected:%v, got:%v",
				i, expected[i], actual[i])
		}
	}
}

// Test modulus under the group
func TestModP(t *testing.T) {
	p := []*large.Int{large.NewInt(17), large.NewIntFromString("717190887961", 10),
		large.NewIntFromString("717190905917", 10)}
	g := large.NewInt(13)

	group := make([]*Group, 0, len(p))
	for i := 0; i < len(p); i++ {
		group = append(group, NewGroup(p[i], g))
	}

	expected := []*large.Int{large.NewInt(2), large.NewIntFromString("269673339004", 10),
		large.NewIntFromString("623940771224", 10)}
	a := []*large.Int{large.NewInt(5000), large.NewIntFromString("beefbeecafe80386", 16),
		large.NewIntFromString("77777777777777777777", 16)}
	actual := make([]*Int, len(expected))
	for i := 0; i < len(expected); i++ {
		actual[i] = group[i].NewInt(1)
		group[i].ModP(a[i], actual[i])
	}

	for i := 0; i < len(expected); i++ {
		if actual[i].value.Cmp(expected[i]) != 0 {
			t.Errorf("TestModP failed, expected: '%v', got: '%v'",
				expected[i].Text(10), actual[i].value.Text(10))
		}
	}

}

// Test that inside panics if cyclicInt doesn't belong to the group
func TestModP_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := large.NewInt(20)
	b := group2.NewInt(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ModP should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	group.ModP(a, b)
}

// Test Inverse under the group
func TestInverse(t *testing.T) {
	p := large.NewInt(17)
	g := large.NewInt(13)
	group := NewGroup(p, g)
	x := group.NewInt(13) //message
	a := group.NewInt(10) //encryption key
	inv := group.NewInt(1)
	inv = group.Inverse(a, inv)             //decryption key
	a = group.Mul(x, a, a)                  // encrypted message
	c := group.Mul(inv, a, group.NewInt(1)) //decrypted message (x)

	if c.value.Cmp(x.value) != 0 {
		t.Errorf("TestInverse failed, expected: '%v', got: '%v'",
			x.value.Text(10), c.value.Text(10))
	}
}

// Test that inverse panics if cyclicInt doesn't belong to the group
func TestInverse_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group.NewInt(20)
	b := group2.NewInt(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Inverse should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	group.Inverse(a, b)
}

// Test Random multiple times to check if
// the number generated is ever outside the group
func TestRandom(t *testing.T) {
	p := large.NewInt(107)
	g := large.NewInt(4)
	group := NewGroup(p, g)
	for i := 0; i < 100000; i++ {
		if !group.Inside(group.Random(group.NewInt(1)).GetLargeInt()) {
			t.Errorf("Generated number is not inside the group!")
		}
	}
}

// Test that Random panics if cyclicInt doesn't belong to the group
func TestRandom_Panic(t *testing.T) {
	p := large.NewInt(107)
	g := large.NewInt(4)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group2.NewInt(20)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Random should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	group.Random(a)
}

type AlwaysErrorReader struct{}

func (r AlwaysErrorReader) Read(b []byte) (int, error) {
	return 1, errors.New("testing reader error")
}

func (r AlwaysErrorReader) SetSeed(seed []byte) error {
	return nil
}

// This test forces random to panic by overwriting the CSPRNG inside the group
func TestRandom_PanicReadErr(t *testing.T) {
	p := large.NewInt(107)
	g := large.NewInt(4)
	group := NewGroup(p, g)

	// Overwrite CSPRNG
	group.rng = AlwaysErrorReader{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Random should panic on read error")
		}
	}()

	group.Random(group.NewInt(1))
}

func TestGen(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)

	// setup array to keep track of frequency of random values
	r := group.NewInt(1)
	rng := make([]int, int(p.Int64()))

	tests := 500
	thresh := 0.3

	// generate randoms
	for i := 0; i < tests; i++ {
		rng[int(group.Random(r).value.Int64())]++
	}

	// make sure 0 and 1 were not generated
	if rng[0] > 0 {
		t.Errorf("TestGen() failed, 0 is outside of the required range")
	}
	if rng[1] > 0 {
		t.Errorf("TestGen() failed, 1 is outside of the required range")
	}

	// check that frequency doesn't exceed threshold
	for i := 0; i < len(rng); i++ {
		if float64(rng[i])/float64(tests) > thresh {
			t.Errorf("TestGen() failed, insufficiently random, value: %v"+
				" occured: %v out of %v tests", i, rng[i], tests)
		}
	}

}

// Test prime getter from the group
func TestGetP(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetP()

	if actual.Cmp(p) != 0 {
		t.Errorf("TestGetP failed, expected: '%v', got: '%v'",
			p.Text(10), actual.Text(10))
	}
}

//Tests the prime byte getter from the group
//Tests the old manual call against the implementation (p.Bytes() vs grp.GetPBytes())
//P is never nil, so no edge case there
func TestGetPBytes(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetPBytes()

	if bytes.Compare(p.Bytes(), actual) != 0 {
		t.Errorf("TestGetPBytes failed, expected: '%v', got: '%v'",
			p.Text(10), actual)
	}

}

// Test generator getter from the group
func TestGetG(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetG()

	if actual.Cmp(g) != 0 {
		t.Errorf("TestGetP failed, expected: '%v', got: '%v'",
			g.Text(10), actual.Text(10))
	}
}

// Test generator getter from the group cyclic version
func TestGetGCyclic(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(33)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetGCyclic()

	if actual.value.Cmp(g) != 0 {
		t.Errorf("TestGetGCyclic failed, expected: '%v', got: '%v'",
			g.Text(10), actual.value.Text(10))
	}
}

// Test prime-1 getter from the group
func TestGetPSub1(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetPSub1()
	ps1 := large.NewInt(16)

	if actual.value.Cmp(ps1) != 0 {
		t.Errorf("TestGetPSub1 failed, expected: '%v', got: '%v'",
			ps1.Text(10), actual.Text(10))
	}
}

// Test prime-1 getter from the group cyclic version
func TestGetPSub1Cyclic(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetPSub1Cyclic()
	ps1 := large.NewInt(16)

	if actual.value.Cmp(ps1) != 0 {
		t.Errorf("TestGetPSub1Cyclic failed, expected: '%v', got: '%v'",
			ps1.Text(10), actual.value.Text(10))
	}

}

// Test (prime-1)/2 getter from the group
func TestGetPSub1Factor(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetPSub1Factor()
	pfactor := large.NewInt(8)

	if actual.Cmp(pfactor) != 0 {
		t.Errorf("TestGetPSub1Factor failed, expected: '%v', got: '%v'",
			pfactor.Text(10), actual.Text(10))
	}
}

// Test (prime-1)/2 getter from the group cyclic version
func TestGetPSub1FactorCyclic(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)
	actual := group.GetPSub1FactorCyclic()
	pfactor := large.NewInt(8)

	if actual.value.Cmp(pfactor) != 0 {
		t.Errorf("TestGetPSub1FactorCyclic failed, expected: '%v', got: '%v'",
			pfactor.Text(10), actual.value.Text(10))
	}
}

// Test array multiplication under the group
func TestArrayMul(t *testing.T) {
	p := large.NewInt(11)
	g := large.NewInt(7)
	grp := NewGroup(p, g)

	expected := large.NewInt(10)

	slc := []*Int{
		grp.NewInt(2),
		grp.NewInt(3),
		grp.NewInt(4),
		grp.NewInt(5),
	}

	c := grp.NewInt(1)
	actual := grp.MulMulti(c, slc...)

	if actual.value.Cmp(expected) != 0 {
		t.Errorf("TestArrayMul failed, expected: '%v', got: '%v'",
			expected, actual)
	}

}

// Test that ArrayMult panics if cyclicInt doesn't belong to the group
func TestArrayMult_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	slc := []*Int{
		group.NewInt(2),
		group2.NewInt(3),
		group.NewInt(4),
		group.NewInt(5),
	}
	a := group.NewInt(20)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ArrayMult should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	group.MulMulti(a, slc...)
}

// Test exponentiation under the group
func TestExp(t *testing.T) {
	p := large.NewInt(117)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	type testStructure struct {
		x *Int
		y *Int
		z *Int
	}

	testStrings := [][]string{
		{"42", "41", "9"},
		{"42", "63", "27"},
		{"69", "42", "27"},
		{"99", "81", "99"},
	}

	var testStructs []testStructure

	for i, strs := range testStrings {
		var ts testStructure

		ts.x = grp.NewIntFromString(strs[0], 10)

		if ts.x == nil {
			t.Errorf("Setup for Test of Exp() for Group failed at 'x' phase of index: %v", i)
		}

		ts.y = grp.NewIntFromString(strs[1], 10)

		if ts.y == nil {
			t.Errorf("Setup for Test of Exp() for Group failed at 'y' phase of index: %v", i)
		}

		ts.z = grp.NewIntFromString(strs[2], 10)

		if ts.z == nil {
			t.Errorf("Setup for Test of Exp() for Group failed at 'z' phase of index: %v", i)
		}

		testStructs = append(testStructs, ts)
	}

	tests := len(testStructs)
	pass := 0

	expected := 0

	for i, testi := range testStructs {
		actual := grp.NewInt(1)
		actual = grp.Exp(testi.x, testi.y, actual)

		result := actual.value.Cmp(testi.z.value)

		if result != expected {
			t.Errorf("Test of Exp() for Group failed at index: %v Expected: %v, %v; Actual: %v, %v",
				i, expected, testi.z.value.Text(10), result, actual.value.Text(10))
		} else {
			pass += 1
		}
	}
	println("Exp() for Group", pass, "out of", tests, "tests passed.")

}

// Test exponentiation of the generator in the group
func TestExpG(t *testing.T) {
	p := large.NewInt(117)
	g := large.NewInt(5)
	grp := NewGroup(p, g)

	type testStructure struct {
		y *Int
		z *Int
	}

	testStrings := [][]string{
		{"42", "64"},
		{"69", "44"},
		{"43", "86"},
		{"2", "25"},
	}

	var testStructs []testStructure

	for i, strs := range testStrings {
		var ts testStructure

		ts.y = grp.NewIntFromString(strs[0], 10)

		if ts.y == nil {
			t.Errorf("Setup for Test of Exp() for Group failed at 'y' phase of index: %v", i)
		}

		ts.z = grp.NewIntFromString(strs[1], 10)

		if ts.z == nil {
			t.Errorf("Setup for Test of Exp() for Group failed at 'z' phase of index: %v", i)
		}

		testStructs = append(testStructs, ts)
	}

	expected := 0

	for i, testi := range testStructs {
		actual := grp.NewInt(1)
		actual = grp.ExpG(testi.y, actual)

		result := actual.value.Cmp(testi.z.value)

		if result != expected {
			t.Errorf("Test of Exp() for Group failed at index: %v Expected: %v; Actual: %v",
				i, testi.z.value.Text(10), actual.value.Text(10))
		}
	}

}

// Test that Exp panics if cyclicInt doesn't belong to the group
func TestExp_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group2.NewInt(20)
	b := group.NewInt(11)
	c := group.NewInt(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Exp should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	group.Exp(a, b, c)
}

// Test random Coprime number generation under the group
func TestRandomCoprime(t *testing.T) {
	// setup test group and generator
	p := large.NewInt(17)
	g := large.NewInt(29)
	group := NewGroup(p, g)

	// setup array to keep track of frequency of random values
	r := group.NewInt(1)
	rng := make([]int, int(p.Int64()))

	tests := 500

	thresh := 0.3

	// generate randoms
	for i := 0; i < tests; i++ {
		rng[int(group.RandomCoprime(r).value.Int64())]++
	}

	// make sure 0 and 1 were not generated
	if rng[0] > 0 {
		t.Errorf("TestRandomeCoprime() failed, 0 is outside of the required range")
	}
	if rng[1] > 0 {
		t.Errorf("TestRandomeCoprime() failed, 1 is outside of the required range")
	}

	// check that frequency doesn't exceed threshold
	for i := 0; i < len(rng); i++ {
		if float64(rng[i])/float64(tests) > thresh {
			t.Errorf("TestRandomCoprime() failed, insufficiently random, value: %v"+
				" occured: %v out of %v tests", i, rng[i], tests)
		}
	}
}

// Test that RandomCoprime panics if cyclicInt doesn't belong to the group
func TestRandomCoprime_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group2.NewInt(20)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("RandomCoprime should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	group.RandomCoprime(a)
}

// This test forces RandomCoprime to panic by overwriting the CSPRNG inside the group
func TestRandomCoprime_PanicReadErr(t *testing.T) {
	p := large.NewInt(5)
	g := large.NewInt(4)
	group := NewGroup(p, g)

	// Overwrite CSPRNG
	group.rng = AlwaysErrorReader{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("RandomCoprime should panic on read error")
		}
	}()

	group.RandomCoprime(group.NewInt(1))
}

//Tests whether the z value is prematurely overwritten in the group after RootCoprime is run
//Tests z value after rootcoprime against a value it's known to be given this input
func TestRootCoprime_ZVal(t *testing.T) {
	p := large.NewInt(17)
	g := large.NewInt(29)

	group := NewGroup(p, g)

	x := group.NewInt(12)
	z := group.NewInt(12)

	//tmp in rootCoPrime is 1 in this case, so z=x**1
	group.RootCoprime(x, z, z)
	if z.value.Cmp(x.value) != 0 {
		t.Errorf("RootCoprime overwrote z value")
	}

}

// You pass a value x = a^y to the RootCoprime function, where y is (supposed to be) coprime with (p-1).
// If y is coprime, then the function returns the value of a
func TestRootCoprime(t *testing.T) {
	tests := 2
	pass := 0

	p := large.NewInt(17)
	g := large.NewInt(29)

	group := NewGroup(p, g)

	a := []*Int{group.NewInt(5), group.NewInt(4), group.NewInt(15)}
	x := group.NewInt(1)
	y := []*Int{group.NewInt(5), group.NewInt(11), group.NewInt(2)}
	z := []*Int{group.NewInt(1), group.NewInt(1), group.NewInt(1)}

	passing := []bool{true, true, false}

	for i := 0; i < 2; i++ {
		group.Exp(a[i], y[i], x)

		tmp := group.RootCoprime(x, y[i], z[i])
		if tmp.value.Cmp(a[i].value) != 0 && passing[i] { //this is wrong, why would z change?
			t.Errorf("RootCoprime Error: Function did not output expected value! a: %v z: %v", a[i].value.Text(16), z[i].value.Text(16))
		} else {
			pass++
		}

	}

	println("RootCoprime", pass, "out of", tests, "tests passed.")
}

// Test that RootCoprime panics if cyclicInt doesn't belong to the group
func TestRootCoprime_Panic(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group.NewInt(20)
	b := group.NewInt(11)
	c := group2.NewInt(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("RootCoprime should panic when one of involved " +
				"cyclicInts doesn't belong to the group")
		}
	}()

	group.RootCoprime(a, b, c)
}

// Test finding a small coprime inverse number in the group
func TestFindSmallCoprimeInverse(t *testing.T) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AACAA68FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	group := NewGroup(p, g)

	num := 1000

	totalBitLen := 0

	bits := uint32(256)

	for i := 0; i < num; i++ {
		z := group.NewInt(1)

		base := group.Random(group.NewInt(1))

		group.FindSmallCoprimeInverse(z, bits)

		zinv := large.NewInt(1).ModInverse(z.value, group.psub1)

		totalBitLen += len(zinv.Bytes()) * 8

		if len(zinv.Bytes())*8 > int(bits) {
			t.Errorf("FindSmallExponent Error: Inverse too large on "+
				"attempt %v; Expected: <%v, Recieved: %v", i, bits,
				uint32(len(zinv.Bytes())*8))
		}

		baseZ := group.NewInt(1)

		group.Exp(base, z, baseZ)

		basecalc := group.NewInt(1)

		basecalc = group.RootCoprime(baseZ, z, basecalc)

		if base.value.Cmp(basecalc.value) != 0 {
			t.Errorf("FindSmallExponent Error: Result incorrect"+
				" on attempt %v; Expected: %s, Recieved: %s", i, base.value.Text(10),
				basecalc.value.Text(10))
		}
	}

	avgBitLen := float32(totalBitLen) / float32(num)

	if float32(avgBitLen) < 0.98*float32(bits) {
		t.Errorf("FindSmallExponent Error: Inverses are not the correct length on average "+
			"; Expected: ~%v, Recieved: %v", 0.95*float32(bits), avgBitLen)
	}

}

// Test finding a small coprime inverse in a group with small p
// This will hit the case where the generated number equals (p-1)/2
func TestFindSmallCoprimeInverse_SmallGroup(t *testing.T) {
	p := large.NewInt(107)
	g := large.NewInt(2)

	group := NewGroup(p, g)
	one := large.NewInt(1)
	num := 1000

	bits := uint32(p.BitLen() - 1)

	for i := 0; i < num; i++ {
		z := group.NewInt(1)

		base := group.Random(group.NewInt(1))

		// z will be unchanged if a number with no inverse is returned
		for z.value.Cmp(one) == 0 {
			group.FindSmallCoprimeInverse(z, bits)
		}

		zinv := large.NewInt(1).ModInverse(z.value, group.psub1)

		if zinv.BitLen() > int(bits) {
			t.Errorf("FindSmallExponent Error: Inverse too large on "+
				"attempt %v; Expected: <%v, Recieved: %v", i, bits,
				zinv.BitLen())
		}

		baseZ := group.NewInt(1)

		group.Exp(base, z, baseZ)

		basecalc := group.NewInt(1)

		basecalc = group.RootCoprime(baseZ, z, basecalc)

		if base.value.Cmp(basecalc.value) != 0 {
			t.Errorf("FindSmallExponent Error: Result incorrect"+
				" on attempt %v; Expected: %s, Recieved: %s", i, base.value.Text(10),
				basecalc.value.Text(10))
		}
	}
}

// Test finding a small coprime inverse in an unsafe group, meaning
// that some numbers don't have an inverse
func TestFindSmallCoprimeInverse_UnsafeGroup(t *testing.T) {
	p := large.NewInt(101)
	g := large.NewInt(2)

	group := NewGroup(p, g)
	one := large.NewInt(1)
	num := 1000

	bits := uint32(6)

	for i := 0; i < num; i++ {
		z := group.NewInt(1)

		base := group.Random(group.NewInt(1))

		// z will be unchanged if a number with no inverse is returned
		for z.value.Cmp(one) == 0 {
			group.FindSmallCoprimeInverse(z, bits)
		}

		zinv := large.NewInt(1).ModInverse(z.value, group.psub1)

		if zinv.BitLen() > int(bits) {
			t.Errorf("FindSmallExponent Error: Inverse too large on "+
				"attempt %v; Expected: <%v, Recieved: %v", i, bits,
				zinv.BitLen())
		}

		baseZ := group.NewInt(1)

		group.Exp(base, z, baseZ)

		basecalc := group.NewInt(1)

		basecalc = group.RootCoprime(baseZ, z, basecalc)

		if base.value.Cmp(basecalc.value) != 0 {
			t.Errorf("FindSmallExponent Error: Result incorrect"+
				" on attempt %v; Expected: %s, Recieved: %s", i, base.value.Text(10),
				basecalc.value.Text(10))
		}
	}
}

// Test that FindSmallCoprimeInverse panics when number of bits is >= log2(p)
func TestFindSmallCoprimeInverse_Panic(t *testing.T) {
	p := large.NewInt(107)
	g := large.NewInt(2)

	group := NewGroup(p, g)
	z := group.NewInt(1)

	bits := uint32(7)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FindSmallCoprimeInverse should panic on bits >= log2(g.prime)")
		}
	}()

	group.FindSmallCoprimeInverse(z, bits)
}

// Test that FindSmallCoprimeInverse panics if cyclicInt doesn't belong to the group
func TestFindSmallCoprimeInverse_PanicArgs(t *testing.T) {
	prime := int64(107)
	p := large.NewInt(prime)
	g := large.NewInt(5)
	group := NewGroup(p, g)
	g2 := large.NewInt(2)
	group2 := NewGroup(p, g2)

	a := group2.NewInt(20)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("RootCoprime should panic when " +
				"cyclicInt doesn't belong to the group")
		}
	}()

	group.FindSmallCoprimeInverse(a, uint32(p.BitLen()-1))
}

// This test forces FindSmallCoprimeInverse to panic by overwriting the CSPRNG inside the group
func TestFindSmallCoprimeInverse_PanicReadErr(t *testing.T) {
	p := large.NewInt(107)
	g := large.NewInt(2)

	group := NewGroup(p, g)

	bits := uint32(p.BitLen() - 1)

	// Overwrite CSPRNG
	group.rng = AlwaysErrorReader{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FindSmallCoprimeInverse should panic on read error")
		}
	}()

	group.FindSmallCoprimeInverse(group.NewInt(1), bits)
}

// Tests that a Group structure that is encoded and then decoded, as a
// gob has the same values.
func TestGroup_GobEncode_GobDecode(t *testing.T) {

	prime := large.NewInt(1000000010101111111)
	gen := large.NewInt(5)
	grp1 := NewGroup(prime, gen)

	b, _ := grp1.GobEncode()

	grp2 := Group{}
	_ = grp2.GobDecode(b)

	if !reflect.DeepEqual(*grp1, grp2) {
		t.Errorf("GobDecode() did not produce the the same original undecoded data\n\treceived: %v\n\texpected: %v", grp1, grp2)
	}
}

// Tests that a Group structure can be marshaled to JSON and unmarshaled to recreate equivalent group
func TestGroup_MarshalJSON_IsValid(t *testing.T) {

	prime := large.NewInt(1000000010101111111)
	gen := large.NewInt(5)
	grp1 := NewGroup(prime, gen)

	// Marshall to bytes
	b, err := grp1.MarshalJSON()

	if err != nil {
		t.Errorf("MarshalJSON() failed to serialize the group: %v", grp1)
	}

	// Unmarshal from bytes
	grp2 := Group{}
	err = grp2.UnmarshalJSON(b)

	if err != nil {
		t.Errorf("UnmarshalJSON() failed to serialize the group: %v", grp1)
	}

	if !reflect.DeepEqual(*grp1, grp2) {
		t.Errorf("UnmarshalJSON() did not produce the the same original undecoded data\n\treceived: %v\n\texpected: %v", grp1, grp2)
	}
}

// BENCHMARKS

func BenchmarkExpForGroup2k(b *testing.B) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AACAA68FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	grp := NewGroup(p, g)

	//prebake inputs
	z := grp.NewInt(1)
	G := grp.GetGCyclic()

	var inputs []*Int
	var outputs []*Int

	for i := 0; i < b.N; i++ {
		nint := grp.Random(grp.NewInt(1))
		inputs = append(inputs, nint)
		outputs = append(outputs, grp.NewInt(1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grp.Exp(G, inputs[i], z)
	}
}

func BenchmarkExpForGroup4k(b *testing.B) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A92108011A723C12A787E6D7" +
		"88719A10BDBA5B2699C327186AF4E23C1A946834B6150BDA" +
		"2583E9CA2AD44CE8DBBBC2DB04DE8EF92E8EFC141FBECAA6" +
		"287C59474E6BC05D99B2964FA090C3A2233BA186515BE7ED" +
		"1F612970CEE2D7AFB81BDD762170481CD0069127D5B05AA9" +
		"93B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934063199" +
		"FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	grp := NewGroup(p, g)

	//prebake inputs
	z := grp.NewInt(1)
	G := grp.GetGCyclic()

	var inputs []*Int
	var outputs []*Int

	for i := 0; i < b.N; i++ {
		nint := grp.Random(grp.NewInt(1))
		inputs = append(inputs, nint)
		outputs = append(outputs, grp.NewInt(1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grp.Exp(G, inputs[i], z)
	}
}

func BenchmarkMulForGroup2k(b *testing.B) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AACAA68FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	grp := NewGroup(p, g)

	//prebake inputs
	z := grp.NewInt(1)

	var inputA []*Int
	var inputB []*Int
	var outputs []*Int

	for i := 0; i < b.N; i++ {
		nint := grp.Random(grp.NewInt(1))
		inputA = append(inputA, nint)
		mint := grp.Random(grp.NewInt(1))
		inputB = append(inputB, mint)
		outputs = append(outputs, grp.NewInt(1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grp.Mul(inputA[i], inputB[i], z)
	}
}

func BenchmarkMulForGroup4k(b *testing.B) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A92108011A723C12A787E6D7" +
		"88719A10BDBA5B2699C327186AF4E23C1A946834B6150BDA" +
		"2583E9CA2AD44CE8DBBBC2DB04DE8EF92E8EFC141FBECAA6" +
		"287C59474E6BC05D99B2964FA090C3A2233BA186515BE7ED" +
		"1F612970CEE2D7AFB81BDD762170481CD0069127D5B05AA9" +
		"93B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934063199" +
		"FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	grp := NewGroup(p, g)

	//prebake inputs
	z := grp.NewInt(1)

	var inputA []*Int
	var inputB []*Int
	var outputs []*Int

	for i := 0; i < b.N; i++ {
		nint := grp.Random(grp.NewInt(1))
		inputA = append(inputA, nint)
		mint := grp.Random(grp.NewInt(1))
		inputB = append(inputB, mint)
		outputs = append(outputs, grp.NewInt(1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grp.Mul(inputA[i], inputB[i], z)
	}
}

func BenchmarkInverse2k(b *testing.B) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AACAA68FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	grp := NewGroup(p, g)

	//prebake inputs
	z := grp.NewInt(1)
	G := grp.GetGCyclic()

	var inputs []*Int
	var outputs []*Int

	for i := 0; i < b.N; i++ {
		nint := grp.Random(grp.NewInt(1))
		nint = grp.Exp(G, nint, z)
		inputs = append(inputs, nint)
		outputs = append(outputs, grp.NewInt(1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grp.Inverse(inputs[i], outputs[i])
	}
}

func BenchmarkInverse4k(b *testing.B) {
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A92108011A723C12A787E6D7" +
		"88719A10BDBA5B2699C327186AF4E23C1A946834B6150BDA" +
		"2583E9CA2AD44CE8DBBBC2DB04DE8EF92E8EFC141FBECAA6" +
		"287C59474E6BC05D99B2964FA090C3A2233BA186515BE7ED" +
		"1F612970CEE2D7AFB81BDD762170481CD0069127D5B05AA9" +
		"93B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934063199" +
		"FFFFFFFFFFFFFFFF"

	p := large.NewIntFromString(primeString, 16)
	g := large.NewInt(2)
	grp := NewGroup(p, g)

	//prebake inputs
	z := grp.NewInt(1)
	G := grp.GetGCyclic()

	var inputs []*Int
	var outputs []*Int

	for i := 0; i < b.N; i++ {
		nint := grp.Random(grp.NewInt(1))
		nint = grp.Exp(G, nint, z)
		inputs = append(inputs, nint)
		outputs = append(outputs, grp.NewInt(1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grp.Inverse(inputs[i], outputs[i])
	}
}

func TestGroup_BytesInside(t *testing.T) {
	test1 := []byte{5}
	test2 := []byte{79}
	test3 := []byte{17}
	test4 := []byte{36}
	grp := NewGroup(large.NewInt(107), large.NewInt(4))

	if !grp.BytesInside(test1, test2, test3, test4) {
		t.Errorf("BytesInside Failed!")
	}
}
