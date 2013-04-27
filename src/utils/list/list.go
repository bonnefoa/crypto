package list

import "math"
import "utils/number"

func Sublist(lst []byte, lenght int) [][]byte {
        size := int(math.Ceil(float64(len(lst)) / float64(lenght)))
        res := make([][]byte, size)
        for i := 0; i < size; i++ {
                index := i*lenght
                indexRight := number.Min(index + lenght, len(lst))
                res[i] = lst[index:indexRight]
        }
        return res
}
