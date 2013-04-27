package hexa

import "strings"

const hex = "0123456789abcdef"

func HexaToBytes(src string) []byte {
        res := make([]byte, len(src) / 2)
        for i := range res {
                first := byte(strings.IndexRune(hex, rune(src[i*2])))
                second := byte(strings.IndexRune(hex, rune(src[i*2 + 1])))
                res[i] = first << 4 | second
        }
        return res
}

func HexaToString(src []byte) string {
        res := make([]uint8, len(src)*2)
        for i, bt := range src {
                first := bt >> 4
                second := bt & 0x0F
                res[i*2] = hex[first]
                res[i*2 + 1] = hex[second]
        }
        return string(res)
}

