package internal

import (
	"math/big"
	"strings"

	"github.com/gofrs/uuid"
)

// NewIFCGUID generate UUID and convert to ifc GUID variant
// Return []byte of ifc GUID
func NewIFCGUID() ([]byte, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return make([]byte, 0), err
	}

	return ByteCompress(u.Bytes()), nil
}

// ByteCompress compress UUID to ifc GUID variant
// Get []byte of UUID without '-', return compresseb []byte
// Random working algorithm, not tested for speed
func ByteCompress(b []byte) []byte {
	var mbi, mbi2, msk big.Int
	var m *big.Int

	guid := make([]byte, 22)
	m = mbi.SetBytes(b)
	mask := msk.SetInt64(63)
	mbi2.SetBytes(b)

	j := 21
	for i := 0; i < 128; i += 6 {
		mask.Lsh(msk.SetInt64(63), uint(i))
		m.And(&mbi2, mask)
		m.Rsh(m, uint(i))
		guid[j] = b64Dict[m.Int64()]
		j--
	}

	return guid
}

// UUIDCompress compress UUID to ifc GUID variant
// Get string representation of UUID without '-', return compresseb []byte
// Random working algorithm, not tested for speed
func UUIDCompress(s string) ([]byte, error) {
	u, err := uuid.FromString(s)

	if err != nil {
		return make([]byte, 0), err
	}

	c := ByteCompress(u.Bytes())
	return c, nil
}

// ByteDecompress decompress ifc GUID to UUID
// Get []byte of ifc GUID
// Random working algorithm, not tested for speed
func ByteDecompress(b []byte) []byte {
	var mbi, msk big.Int
	uuid := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	mbi.FillBytes(uuid)
	d := 0
	for i := 0; i < len(b); i++ {
		d = int(strings.Index(b64Dict, string(b[i])))
		mbi.Add(&mbi, msk.SetInt64(int64(d)))
		mbi.Lsh(&mbi, 6)
	}
	mbi.Rsh(&mbi, 6)
	j := 15
	for i := len(mbi.Bytes()) - 1; i > -1; i-- {
		uuid[j] = mbi.Bytes()[i]
		j--
	}
	return uuid
}
