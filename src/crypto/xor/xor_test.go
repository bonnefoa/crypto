package xor

import "testing"
import "testing/assert"
import "encoding/hexa"
import "encoding/base64"
import "io/ioutil"

func TestXor(t *testing.T) {
        src := hexa.HexaToBytes("1c0111001f010100061a024b53535009181c")
        xor := hexa.HexaToBytes("686974207468652062756c6c277320657965")
        expected := hexa.HexaToBytes("746865206b696420646f6e277420706c6179")
        assert.AssertBytesEquals(t, Xor(src, xor), expected)
}

func TestSeekXor(t *testing.T) {
        src := hexa.HexaToBytes("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
        phrases := SeekXor(src)
        assert.AssertEquals(t, string(phrases[0].Decrypted), "Cooking MC's like a pound of bacon")
}

func TestSeekXorFromFile(t *testing.T) {
        filename := "./source_4"
        phrases := SeekXorInFile(filename)
        assert.AssertEquals(t, string(phrases[0].Decrypted), "Now that the party is jumping\n")
}

func TestRepeatedXor(t *testing.T) {
        src := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
        xor := "ICE"
        res := Xor([]byte(src), []byte(xor))
        expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
        assert.AssertEquals(t, hexa.HexaToString(res), expected)
}

func TestDistance(t *testing.T) {
        first := []byte("this is a test")
        snd := []byte("wokka wokka!!!")
        assert.AssertEquals(t, HammingDistance(first, snd), 37)
}

func TestBreakCypher(t *testing.T) {
        content, _ := ioutil.ReadFile("source_6")
        input := base64.StringBase64ToBytes(string(content))
        guesses := GuessProbableKeysize(2, 40, input)
        bests := guesses[:5]
        for _, guess := range bests {
                candidates := GetHistograms(guess.size, content)
                for _, candidate := range candidates {
                        t.Logf("res %q\n", candidate)
                }
        }
        t.Fail()
}
