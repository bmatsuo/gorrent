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
    dsumm := func(s string, d interface{}) string { return fmt.Sprintf("%s->%d", s, d) }

    d := NewDecoder([]byte(in))
    i, err := d.Decode()
    if !exp_err {
        if err != nil {
            DecodingError(t, "int", "unexpected error", dsumm(in, exp), err.String())
        }
        if i != exp {
            DecodingError(t, "int", "unexpected result", strconv.Itoa64(exp), dsumm(in, i))
        }
    } else {
        if err == nil {
            DecodingError(t, "int", "unexpected result", "Error", dsumm(in, i))
        }
    }
}

func TestInteger(t *testing.T) {
    it(t, "i23e", 23, false)
    it(t, "i124145124e", 124145124, false)
    it(t, "i15155", 0, true)
    it(t, "55", 55, true)
}

func st(t *testing.T, in string, exp string, exp_err bool) {
    // Summarize a decoding, either expected or observed.
    dsumm := func(s string, d interface{}) string { return fmt.Sprintf("%s->%s", s, d) }

    d := NewDecoder([]byte(in))
    s, err := d.Decode()
    if !exp_err {
        if err != nil {
            DecodingError(t, "string", "unexpected error", dsumm(in, exp), err.String())
        }
        if s != exp {
            DecodingError(t, "string", "unexpected result", exp, dsumm(in, s))
        }
    } else {
        if err == nil {
            DecodingError(t, "string", "unexpected result", "Error", dsumm(in, s))
        }
    }
}

func TestString(t *testing.T) {
    st(t, "5:hello", "hello", false)
    st(t, "6:world", "world", true)
}

func lt(t *testing.T, in string, exp []interface{}, exp_err bool) {
    // Summarize a decoding, either expected or observed.
    dsumm := func(s string, list []interface{}) string { return fmt.Sprintf("%s->%v", s, list) }

    d := NewDecoder([]byte(in))
    l, err := d.Decode()
    if !exp_err {
        if err != nil {
            DecodingError(t, "list", "unexpected error", dsumm(in, exp), err.String())
        }
        switch l.(type) {
        case nil:
            if len(exp) != 0 {
                DecodingError(t, "list", "unexpected result", fmt.Sprintf("%v", exp), "nil")
            }
        case []interface{}:
            list := l.([]interface{})
            if len(list) != len(exp) {
                DecodingError(t, "list", "unexpected result", fmt.Sprintf("%v", exp), dsumm(in, list))
            } else {
                for i := range list {
                    if list[i] != exp[i] {
                        DecodingError(t, "list", "unexpected result", fmt.Sprintf("%v", exp), dsumm(in, list))
                        break
                    }
                }
            }
        }
    } else {
        if err == nil {
            list := l.([]interface{})
            DecodingError(t, "string", "unexpected result", "Error", dsumm(in, list))
        }
    }
}

func TestList(t *testing.T) {
    lt(t, "li124145124ee", []interface{}{int64(124145124)}, false)
    lt(t, "li15155ee", []interface{}{int64(15155)}, false)
    lt(t, "le", []interface{}{}, false)
    lt(t, "li15155e", []interface{}{}, true)
}
