package internal

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	b16Dict = "0123456789abcdef"
	b64Dict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_$"
)

// Encode2HexString encode UTF8 string to string in UTF IFC format
func Encode2HexString(in string) (out string, err error) {
	var a string
	fmt.Println(in)
	for _, rv := range in {
		if int(rv) > 126 {
			a = a + `\X2\` + strings.ToUpper(fmt.Sprintf("%04s", strconv.FormatInt(int64(rv), 16))) + `\X0\`
		} else {
			a = a + string(rv)
		}
	}
	return a, nil
}
