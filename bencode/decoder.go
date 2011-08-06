/*
	Package bencode implements reading and writing of 'bencoded'
	object streams used by the Bittorent protocol.

*/
package bencode

import (
    "strconv"
    "fmt"
    "os"
)

//A Decoder reads and decodes bencoded objects from an input stream.
//It returns objects that are either an "Integer", "String", "List" or "Dict".
//
//Example usage:
//	d := bencode.NewDecoder([]byte("i23e4:testi123e"))
//	for !p.Consumed {
//		o, _ := p.Decode()
//		fmt.Printf("obj(%s): %#v\n", reflect.TypeOf(o).Name, o)
//	}
type Decoder struct {
    stream   []byte
    pos      int
    Consumed bool //true if we have consumed all tokens
}

//NewDecoder creates a new decoder for the given token stream
func NewDecoder(b []byte) *Decoder { return &Decoder{b, 0, false} }

//Decode reads one object from the input stream
func (self *Decoder) Decode() (res interface{}, err os.Error) {
    return self.nextObject()
}

var (
    ErrorConsumed = os.NewError("This parser's token stream is consumed!")
    ErrorNoTerminator = os.NewError("No terminating 'e' found!")
)

//DecodeAll reads all objects from the input stream
func (self *Decoder) DecodeAll() (res []interface{}, err os.Error) {
    var obj interface{}
    for err = ErrorConsumed; !self.Consumed; err = nil {
        if obj, err = self.nextObject(); err != nil {
            return
        }
        res = append(res, obj)
    }
    return
}

//fetch the next object at position 'pos' in 'stream'
func (self *Decoder) nextObject() (res interface{}, err os.Error) {
    if self.Consumed {
        return nil, ErrorConsumed
    }

    switch c := self.stream[self.pos]; {
    case c == 'i':
        res, err = self.nextInteger()
    case c >= '0' && c <= '9':
        res, err = self.nextString()
    case c == 'l':
        res, err = self.nextList()
    case c == 'd':
        res, err = self.nextDict()
    default:
        res = nil
        err = fmt.Errorf("Couldn't parse '%s' ... '%s' (%d)", self.stream, string(self.stream[self.pos]), self.pos)
    }
    if self.pos >= len(self.stream) {
        self.Consumed = true
    }
    return
}

//fetches next integer from stream and advances pos pointer
func (self *Decoder) nextInteger() (res int64, err os.Error) {
    if self.stream[self.pos] != 'i' {
        return 0, os.NewError("No starting 'i' found")
    }
    validstart := false //flag to check for leading 0's
    idx := self.pos + 1
    for {
        if self.stream[idx] == 'e' {
            break
        }

        if self.stream[idx] == '0' && !validstart {
            err = os.NewError("Leading Zeros are not allowed in bencoded integers!")
            return
        }

        //check for bytes != '-' and '0'..'9'
        if (self.stream[idx] < '0' || self.stream[idx] > '9') && self.stream[idx] != '-' {
            err = os.NewError("Invalid byte '" + string(self.stream[idx]) + "' in encoded integer.")
            return
        }
        validstart = true

        if idx++; idx >= len(self.stream) {
            return 0, ErrorNoTerminator
        }
    }

    s := string(self.stream[self.pos+1 : idx])
    if res, err = strconv.Atoi64(s); err != nil {
        return // Or: return 0, err
    }
    self.pos = idx + 1

    return
}

//fetches next string from stream and advances pos pointer
func (self *Decoder) nextString() (res string, err os.Error) {
    if self.stream[self.pos] < '0' || self.stream[self.pos] > '9' {
        err = os.NewError("No string length determinator found")
        return
    }

    len_start := self.pos
    len_end := self.pos

    //scan length
    for self.stream[len_end] != ':' {
        if len_end++; len_end >= len(self.stream) {
            err = os.NewError("No string found ...")
            return
        }
    }

    if l, e := strconv.Atoi(string(self.stream[len_start:len_end])); e != nil {
        err = fmt.Errorf("Couldn't parse string length specifier: %s", e.String())
    } else if l >= len(self.stream[len_end:]) {
        err = os.NewError("Specified length longer than data buffer ...")
    } else {
        len_end++ //skip the ':'
        res = string(self.stream[len_end : len_end+l])
        err = nil
        self.pos = len_end + l
    }
    return
}

//fetches a list (and its contents) from stream and advances pos
func (self *Decoder) nextList() (res []interface{}, err os.Error) {
    if self.stream[self.pos] != 'l' {
        err = os.NewError("This is not a list!")
        return
    }
    self.pos++ //skip 'l'

    if self.stream[self.pos] == 'e' {
        self.pos++ //skip 'e'
        return
    }

    var obj interface{}
    for {
        if obj, err = self.nextObject(); err != nil {
            return
        }
        res = append(res, obj)
        if self.pos >= len(self.stream) {
            err = ErrorNoTerminator
            return
        }
        if self.stream[self.pos] == 'e' {
            self.pos++ //skip 'e'
            break
        }
    }
    return
}

//fetches a dict
//bencoded dicts must have their keys sorted lexically. but I guess
//we can ignore that and work with unsorted maps. (wtf?! sorted maps ...)
func (self *Decoder) nextDict() (res map[string]interface{}, err os.Error) {
    if self.stream[self.pos] != 'd' {
        err = os.NewError("This is not a dict!")
        return
    }
    self.pos++ //skip 'd'

    if self.stream[self.pos] == 'e' {
        self.pos++ //skip 'e'
        return
    }

    res = make(map[string]interface{})
    var (
        key string
        val interface{}
    )
    for {
        if key, err = self.nextString(); err != nil {
            return
        }
        if val, err = self.nextObject(); err != nil {
            return
        }
        //fmt.Printf("key: %s\nval: %#v\n", key, val)
        res[string(key)] = val
        if self.pos >= len(self.stream) {
            err = ErrorNoTerminator
            return
        }
        if self.stream[self.pos] == 'e' {
            self.pos++ //skip 'e'
            break
        }
    }
    return
}
