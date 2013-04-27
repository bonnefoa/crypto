package hexa

import "testing"
import "testing/assert"

func TestHexa(t *testing.T) {
        assert.AssertEquals(t, HexaToString([]byte{0x0f}), "0f")
        assert.AssertEquals(t, HexaToString([]byte{0x0f, 0xac}), "0fac")
}

func TestHexaToBytes(t *testing.T) {
        assert.AssertBytesEquals(t, HexaToBytes("0f"), []byte{0x0f})
}
