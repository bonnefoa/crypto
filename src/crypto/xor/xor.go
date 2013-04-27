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

type Phrase struct {
        Source []byte
        Decrypted []byte
        Xor []byte
        Score int
}

func (phr Phrase) String() string {
        return fmt.Sprintf("Decrypted %q, source %q, xor %q, score %d", phr.Decrypted,
                phr.Source, phr.Xor, phr.Score)
}

type Phrases []*Phrase
type ByScore struct { Phrases }

func (s Phrases) Len() int { return len(s) }
func (s Phrases) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByScore) Less(i, j int) bool { return s.Phrases[i].Score > s.Phrases[j].Score }

func Xor(src []byte, xor []byte) []byte {
        res := make([]byte, len(src))
        for i, bt := range src {
                res[i] = bt ^ xor[i % len(xor)]
        }
        return res
}

func ScoreText(src []byte) int {
        score := 0
        for _, r := range(bytes.Runes(src)) {
                if unicode.IsLetter(r) {
                        score++
                }
                if unicode.IsLower(r) {
                        score++
                }
                if unicode.IsSpace(r) {
                        score += 4
                }
                if unicode.IsPunct(r) {
                        score -= 2
                }
                if unicode.IsSymbol(r) {
                        score -= 2
                }
                if !unicode.IsPrint(r) {
                        score -= 20
                }
        }
        return score
}

func SeekXor(src []byte) Phrases {
        phrases := make(Phrases, 0)
        var i uint8
        for i = 0; i < 255; i++ {
                xors := []byte {byte(i)}
                xored := Xor(src, xors)
                phrase := new(Phrase)
                *phrase = Phrase{src, xored, xors, ScoreText(xored)}
                phrases = append(phrases, phrase)
        }
        sort.Sort(ByScore{phrases})
        return phrases
}

func SeekXorInFile(filename string) Phrases {
        content, _ := ioutil.ReadFile(filename)
        lines := strings.Split(string(content), "\n")
        phrases := make(Phrases, 0)
        for _, line := range lines {
                for _, phrase := range SeekXor(hexa.HexaToBytes(line)) {
                        phrases = append(phrases, phrase)
                }
        }
        sort.Sort(ByScore{phrases})
        return phrases
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
        guesses := make(KeysizeGuesses, 0)
        for i := min; i < max; i++ {
                firstBatch := src[0:i]
                secondBatch := src[i:i*2]
                dist := HammingDistance(firstBatch, secondBatch)
                normalized := float64(dist) / float64(i)
                guess := new(KeysizeGuess)
                *guess = KeysizeGuess{i, dist, normalized}
                guesses = append(guesses, guess)
        }
        sort.Sort(guesses)
        return guesses
}


func GetHistograms(size int, src []byte) Phrases {
        log.Printf("Trying key size %d \n", size)
        sublists := list.Sublist(src, size)
        blocks := make([][]byte, size)
        for i:=0; i < size; i++ {
                blocks[i] = make([]byte, len(sublists))
        }
        for i, subblock := range sublists {
                for j, bt := range subblock {
                        blocks[j][i] = bt
                }
        }

        phrases := make(Phrases, size)
        for i, subblock := range blocks {
                candidates := SeekXor(subblock)
                phrases[i] = candidates[0]
                for j:=0; j < 3; j++ {
                        log.Printf("Candidates %d for %d / %d : %s", j, i, size, candidates[j])
                }
        }
        return phrases
}
