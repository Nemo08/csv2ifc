package internal

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestByteCompress(t *testing.T) {
	var testData = map[string]string{
		"0ey9sVfFPDdgqxlMM2IzdH": "28f09d9f-a4f6-4d9e-ad3b-bd65824bd9d1",
		"0F3WqC2me920S61GG30W40": "0f0e0d0c-0b0a-0908-0706-050403020100",
		"0000000000000000000000": "00000000-0000-0000-0000-000000000000",
	}
	var u uuid.UUID
	var b []byte

	for k, v := range testData {
		err := u.UnmarshalText([]byte(v))
		if err != nil {
			t.Error("Expected ", k, " and ", v, " got error ", err)
		}

		mb, err := u.MarshalBinary()
		if err != nil {
			t.Error("Expected ", k, " and ", v, " got error ", err)
		}

		b = ByteCompress(mb)

		if string(b) != k {
			t.Error("Expected ", k, " got ", string(b))
		}
	}
}

func TestByteDecompress(t *testing.T) {
	var testData = map[string]string{
		"0ey9sVfFPDdgqxlMM2IzdH": "28f09d9f-a4f6-4d9e-ad3b-bd65824bd9d1",
		"0F3WqC2me920S61GG30W40": "0f0e0d0c-0b0a-0908-0706-050403020100",
		"0000000000000000000000": "00000000-0000-0000-0000-000000000000",
	}
	var b []byte

	for k, v := range testData {
		b = ByteDecompress([]byte(k))
		u, err := uuid.FromBytes(b)
		if err != nil {
			t.Error("Decompress error for ", k, " got ", err)
		}
		if u.String() != v {
			t.Error("Expected ", v, " got ", u.String())
		}
	}
}

func TestUUIDCompress(t *testing.T) {
	var testData = map[string]string{
		"0ey9sVfFPDdgqxlMM2IzdH": "28f09d9f-a4f6-4d9e-ad3b-bd65824bd9d1",
		"0F3WqC2me920S61GG30W40": "0f0e0d0c-0b0a-0908-0706-050403020100",
		"0000000000000000000000": "00000000-0000-0000-0000-000000000000",
	}

	for k, v := range testData {
		b, err := UUIDCompress(v)
		if err != nil {
			t.Error("Expected ", k, " and ", string(b), " got ", err)
		}
		if string(b) != k {
			t.Error("Expected ", k, " got ", string(b))
		}
	}
}
