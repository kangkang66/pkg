package xaes

import (
	"fmt"
	"testing"
)

func TestNewAes(t *testing.T) {
	aes := NewAes("4ea5a27becaa2e054b950b3b7c8c95034949e4vb")
	val,err := aes.Encrypt([]byte("kangkang"))
	fmt.Println(val,err)

	dstr :="FtEDCQ8kRMeijBZ+CoJzidB5g902MdGdug58E1H+i7s="
	dval,err := aes.Decrypt(dstr)
	fmt.Println(string(dval),err)
}
