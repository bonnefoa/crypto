package xor

import "testing"
import "testing/assert"
import "crypto/hexa"
//import "io/ioutil"
//import "strings"
//import "crypto/base64"

const testText = "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
const testKey = "rxtcfyvg"

func TestXor(t *testing.T) {
        src := hexa.HexaToBytes("1c0111001f010100061a024b53535009181c")
        xor := hexa.HexaToBytes("686974207468652062756c6c277320657965")
        expected := hexa.HexaToBytes("746865206b696420646f6e277420706c6179")
        assert.AssertBytesEquals(t, Xor(src, xor), expected)
}

func TestSeekXor(t *testing.T) {
        src := hexa.HexaToBytes("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
        candidates := SeekXor(src)
        for i:=0; i<5; i++ {
                t.Logf("Candidate %d is %q\n", candidates[i])
        }
        assert.AssertEquals(t, string(candidates[0].Decrypted), "Cooking MC's like a pound of bacon")
}

func TestSeekXorFromFile(t *testing.T) {
        filename := "./source_4"
        candidates := SeekXorInFile(filename)
        for i:=0; i < 5; i++ {
                t.Logf("res %q\n", candidates[i])
        }
        assert.AssertEquals(t, string(candidates[0].Decrypted), "Now that the party is jumping\n")
}

func TestRepeatedXor(t *testing.T) {
        xor := "ICE"
        res := Xor([]byte(testText), []byte(xor))
        expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
        assert.AssertEquals(t, hexa.HexaToString(res), expected)
}

func TestDistance(t *testing.T) {
        first := []byte("this is a test")
        snd := []byte("wokka wokka!!!")
        assert.AssertEquals(t, HammingDistance(first, snd), 37)
}

func TestGuessKeySize(t *testing.T) {
        key := []byte(testKey)
        for i:=2; i < len(key); i++ {
                xoredText := Xor([]byte(testText), key[0:i])
                guesses := GuessProbableKeysize(2, 8, xoredText)
                t.Logf("Guessing %d\n", i)
                for _, guess := range guesses {
                        t.Logf("Guess is %s\n", guess)
                }
                found := false
                for k := 0; k < 3; k++ {
                        if guesses[k].size == i {
                                found = true
                        }
                }
                if ! found { t.FailNow() }
        }
}

func TestBreakCypher(t *testing.T) {
        src := []byte(testText)
        srcKey := []byte(testKey)
        for i:=2; i < len(srcKey); i++ {
                xor := Xor(src, srcKey[0:i])
                guesses := GuessProbableKeysize(2, len(srcKey), xor)
                for _, guess := range guesses[:2] {
                        candidates := GuessRepeatXor(guess.size, xor)
                        found := false
                        for k:=0; k<3; k++ {
                                guessKey := candidates[k].Xor
                                t.Logf("guess %q\n", guessKey)
                                //if guessKey == srcKey { found = true }
                        }
                        if !found {t.FailNow()}
                }
        }
}

//func TestBreakCypher(t *testing.T) {
        //content, _ := ioutil.ReadFile("source_6")
        //stripped := strings.Replace(string(content), "\n", "", 0)
        //input := base64.StringBase64ToBytes(stripped)
        //guesses := GuessProbableKeysize(2, 40, input)
        //bests := guesses[:30]
        //for _, guess := range bests {
                //candidates := GuessRepeatXor(guess.size, input)
                //for _, candidate := range candidates {
                        //t.Logf("res %q\n", candidate)
                //}
        //}
        //t.Fail()
//}
