package base64

import "testing"
import "testing/assert"
import "crypto/hexa"

func TestGroupBySize(t *testing.T) {
        res := groupBySix( []byte("Man") )
        assert.AssertBytesEquals(t, res, []byte { 0x13, 0x16, 0x05, 0x2e } )
}

func TestSimpleBase64(t *testing.T) {
        src := []byte("Man")
        expected := "TWFu"
        assert.AssertEquals(t, Base64ToString(Base64Encode(src)), expected)
}

func TestBase64Encode(t *testing.T) {
        src := hexa.HexaToBytes("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")
        expected := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
        assert.AssertEquals(t, Base64ToString(Base64Encode(src)), expected)
}

func TestStringToBase64(t *testing.T) {
        src := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
        assert.AssertEquals(t, Base64ToString(StringToBase64(src)), src)
}

func TestBase64Decode(t *testing.T) {
        src := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
        assert.AssertEquals(t, Base64ToString(StringToBase64(src)), src)
}
