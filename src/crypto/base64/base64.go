package base64

import "strings"
import "utils/list"

func groupBySix(group []byte) []byte {
        first := group[0] >> 2
        second := (group[0] & 0x03) << 4 | (group[1] >> 4)
        third := (group[1] & 0x0f) << 2 | (group[2] & 0xc0) >> 6
        fourth := group[2] & 0x3f
        return []byte{first, second, third, fourth}
}

func mergeGroupOfSix(group []byte) []byte {
        first := group[0] << 2 | group[1] >> 4
        second := (group[1] & 0x0f) << 4 | (group[2] >> 4)
        third := (group[2] & 0x03) << 6 | group[3]
        return []byte{first, second, third}
}

const base64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func Base64ToString(src []byte) string {
        res := make([]uint8, len(src))
        for i, v := range src {
                res[i] = base64[v]
        }
        return string(res)
}

func Base64Encode(src []byte) []byte {
        sublists := list.Sublist(src, 3)
        res := make([]byte, len(sublists) * 4)
        for i, group := range sublists {
                subgroup := groupBySix(group)
                copy(res[i*4:i*4+4], subgroup)
        }
        return res
}

func Base64Decode(src []byte) []byte {
        sublists := list.Sublist(src, 4)
        res := make([]byte, len(sublists) * 3)
        for i, group := range sublists {
                subgroup := mergeGroupOfSix(group)
                copy(res[i*3:i*3+3], subgroup)
        }
        return res
}

func StringToBase64(src string) []byte {
        res := make([]byte, len(src))
        for i, rn := range src {
                res[i] = byte(strings.IndexRune(base64, rn))
        }
        return res
}

func StringBase64ToBytes(src string) []byte {
        return Base64Decode(StringToBase64(src))
}
