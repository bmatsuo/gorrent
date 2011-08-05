package bencode

import (
	"testing"
    "strconv"
    "fmt"
)

func DecodingError(t *testing.T, typ, msg, exp, recv string) {
    t.Errorf("Decoding %s: %s (expected %s) %s", typ, msg, exp, recv)
}

func it(t *testing.T, in string, exp int64, exp_err bool) {
    // Summarize a decoding, either expected or observed.
    dsumm := func(s string, d interface{}) string { return fmt.Sprintf("%s->%d", s, d.(int64)) }

	d := NewDecoder([]byte(in))
    i, err := d.Decode()
	if !exp_err {
		if err != nil {
            DecodingError(t, "Int", "unexpected error", dsumm(in, exp), err.String())
		}
		if i != exp {
            DecodingError(t, "Int", "unexpected result", strconv.Itoa64(exp), dsumm(in, i))
		}
	} else {
		if err == nil {
            DecodingError(t, "Int", "unexpected result", "Error", dsumm(in, i))
		}
	}
}

func TestInteger(t *testing.T) {
	it(t, "i23e", 23, false)
	it(t, "i124145124e", 124145124, false)
	it(t, "i15155", 0, true)
	it(t, "55", 55, true)
}
