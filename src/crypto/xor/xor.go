package xor

import "unicode"
import "bytes"
import "sort"
import "fmt"
import "io/ioutil"
import "strings"
import "crypto/hexa"
import "utils/list"
import "log"

type XorCandidate struct {
        Source []byte
        Decrypted []byte
        Xor []byte
        Score float32
}

func (x XorCandidate) String() string {
        return fmt.Sprintf("xor %x, score %f, Decrypted %q", x.Xor,
                x.Score, x.Decrypted)
}

type XorCandidates []*XorCandidate
type ByScore struct { XorCandidates }

func (s XorCandidates) Len() int { return len(s) }
func (s XorCandidates) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByScore) Less(i, j int) bool { return s.XorCandidates[i].Score > s.XorCandidates[j].Score }

func Xor(src []byte, xor []byte) []byte {
        res := make([]byte, len(src))
        for i, bt := range src {
                res[i] = bt ^ xor[i % len(xor)]
        }
        return res
}

func ScoreText(src []byte) float32 {
        score := 0
        for _, r := range(bytes.Runes(src)) {
                if unicode.IsLetter(r) {
                        score++
                }
                if unicode.IsLower(r) {
                        score += 4
                }
                if unicode.IsSpace(r) {
                        score += 30
                }
                if unicode.IsPunct(r) {
                        score -= 2
                }
                if unicode.IsSymbol(r) {
                        score -= 30
                }
                if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
                        score -= 20
                }
        }
        return float32(score) / float32(len(src))
}

func SeekXor(src []byte) XorCandidates {
        xorCandidates := make(XorCandidates, 0)
        var i uint8
        for i = 1; i < 255; i++ {
                xors := []byte {byte(i)}
                xored := Xor(src, xors)
                xorCandidate := new(XorCandidate)
                *xorCandidate = XorCandidate{src, xored, xors, ScoreText(xored)}
                xorCandidates = append(xorCandidates, xorCandidate)
        }
        sort.Sort(ByScore{xorCandidates})
        return xorCandidates[0:5]
}

func SeekXorInFile(filename string) XorCandidates {
        content, _ := ioutil.ReadFile(filename)
        lines := strings.Split(string(content), "\n")
        xorCandidates := make(XorCandidates, 0)
        for _, line := range lines {
                for _, xorCandidate := range SeekXor(hexa.HexaToBytes(line)) {
                        xorCandidates = append(xorCandidates, xorCandidate)
                }
        }
        sort.Sort(ByScore{xorCandidates})
        return xorCandidates
}


func AbsInt(num int) int {
        if num < 0 { return -num }
        return num
}

func HammingDistance(first []byte, snd []byte) int {
        distance := 0
        for i := range first {
                val := first[i] ^ snd[i]
                for val != 0 {
                        distance++
                        val = val & (val - 1)
                }
        }
        return distance
}

type KeysizeGuess struct {
        size int
        editDistance int
        normalized float64
}

func (g KeysizeGuess) String() string {
        return fmt.Sprintf("Size %d, editDistance %d, normalizeDist %f", g.size,
                g.editDistance, g.normalized)
}

type KeysizeGuesses []*KeysizeGuess

func (s KeysizeGuesses) Len() int { return len(s) }
func (s KeysizeGuesses) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s KeysizeGuesses) Less(i, j int) bool { return s[i].normalized < s[j].normalized }

func GuessProbableKeysize(min, max int, src []byte) KeysizeGuesses {
        numBatch := 4
        guesses := make(KeysizeGuesses, 0)
        for i := min; i < max; i++ {
                dist := 0
                for j:=0; j < numBatch; j++ {
                        offset := j*i
                        firstBatch := src[offset:offset + i]
                        secondBatch := src[offset + i:offset + i*2]
                        dist += HammingDistance(firstBatch, secondBatch)
                }
                normalized := (float64(dist) / float64(i)) / float64(numBatch)
                guess := new(KeysizeGuess)
                *guess = KeysizeGuess{i, dist, normalized}
                guesses = append(guesses, guess)
        }
        sort.Sort(guesses)
        return guesses
}

func mergeXorCandidates(xorCandidates XorCandidates) *XorCandidate {
        size := len(xorCandidates)
        strSize := len(xorCandidates[0].Source)

        src := make([]byte, size * strSize)
        decrypted := make([]byte, size * strSize)
        xor := make([][]byte, 0)
        score := float32(0)
        for i, xorCandidate := range xorCandidates {
                for j:=0; j < strSize; j++ {
                        src[i + j * size] = xorCandidate.Source[j]
                        decrypted[i + j * size] = xorCandidate.Decrypted[j]
                }
                xor = append(xor, xorCandidate.Xor)
                score += xorCandidate.Score
        }
        xorCandidate := new(XorCandidate)
        *xorCandidate = XorCandidate{src,
                decrypted,
                bytes.Join(xor, []byte{}), score}
        return xorCandidate
}

func TransposeBlocks(src [][]byte) [][]byte {
        n := len(src)
        m := len(src[0])
        blocks := make([][]byte, n)
        for i:=0; i < n; i++ {
                blocks[i] = make([]byte, m)
        }
        for i, subblock := range src {
                for j, bt := range subblock {
                        blocks[j][i] = bt
                }
        }
        return blocks
}

func GuessRepeatXor(size int, src []byte) XorCandidates {
        log.Printf("Trying key size %d \n", size)
        sublists := list.Sublist(src, size)
        blocks := TransposeBlocks(sublists)

        xorCandidates := make(XorCandidates, 0)
        xorBlocks := make([]XorCandidates, size)
        for i, subblock := range blocks {
                candidates := SeekXor(subblock)
                xorBlocks[i] = candidates
        }
        for i:=0; i < 3; i++ {
                batch := make(XorCandidates, 0)
                for j := 0; i < size; j++ {
                        log.Printf("Got %d/%d\n", j, i)
                        batch = append(batch, xorBlocks[j][i])
                }
                xorCandidates = append(xorCandidates,
                        mergeXorCandidates(batch))
        }
        sort.Sort(ByScore{xorCandidates})
        return xorCandidates
}
