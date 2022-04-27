package internal

import (
	"testing"
)

func TestEncode2HexString(t *testing.T) {
	es, err := Encode2HexString("hbbffff ТРОГялпжRREНЕПД")
	if err != nil {
		t.Error("Encode2HexString error", err)
	}
	if es != `hbbffff \X2\0422\X0\\X2\0420\X0\\X2\041e\X0\\X2\0413\X0\\X2\044f\X0\\X2\043b\X0\\X2\043f\X0\\X2\0436\X0\RRE\X2\041d\X0\\X2\0415\X0\\X2\041f\X0\\X2\0414\X0\` {
		t.Error(`Expected hbbffff \X2\0422\X0\\X2\0420\X0\\X2\041e\X0\\X2\0413\X0\\X2\044f\X0\\X2\043b\X0\\X2\043f\X0\\X2\0436\X0\RRE\X2\041d\X0\\X2\0415\X0\\X2\041f\X0\\X2\0414\X0\ got `, es)
	}
}
