/*
	Package bencode implements reading and writing of 'bencoded'
	object streams used by the Bittorent protocol.

*/
package bencode

import (
	"errors"
	"fmt"
	"strconv"
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
func (self *Decoder) Decode() (res interface{}, err error) {
	return self.nextObject()
}

var (
	ErrorConsumed     = errors.New("This parser's token stream is consumed!")
	ErrorNoTerminator = errors.New("No terminating 'e' found!")
)

//DecodeAll reads all objects from the input stream
func (self *Decoder) DecodeAll() (res []interface{}, err error) {
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
func (self *Decoder) nextObject() (res interface{}, err error) {
	if self.Consumed {
		return nil, ErrorConsumed
	}

	switch c := self.stream[self.pos]; c {
	case 'i':
		res, err = self.nextInteger()
	case 'l':
		res, err = self.nextList()
	case 'd':
		res, err = self.nextDict()
	default:
		if c >= '0' && c <= '9' {
			res, err = self.nextString()
		} else {
			err = fmt.Errorf("Couldn't parse '%s' index %d (%s)", self.stream, self.pos, string(self.stream[self.pos]))
		}
	}
	if self.pos >= len(self.stream) {
		self.Consumed = true
	}
	return
}

//fetches next integer from stream and advances pos pointer
func (self *Decoder) nextInteger() (res int64, err error) {
	if self.stream[self.pos] != 'i' {
		return 0, errors.New("No starting 'i' found")
	}
	self.pos++
	idx := self.pos

	if self.stream[idx] == '-' {
		idx++
	}
	start := idx

	for self.stream[idx] != 'e' {
		//check for bytes != '-' and '0'..'9'
		if self.stream[idx] < '0' || self.stream[idx] > '9' {
			err = fmt.Errorf("Invalid byte '%s' in encoded integer.", string(self.stream[idx]))
			return
		}

		if idx++; idx >= len(self.stream) {
			return 0, ErrorNoTerminator
		}
	}

	if start == idx {
		err = errors.New("No bytes in integer")
		return
	}
	if self.stream[start] == '0' && idx-start > 1 {
		err = errors.New("Leading Zeros are not allowed in bencoded integers!")
		return
	}

	s := string(self.stream[self.pos:idx])
	if res, err = strconv.ParseInt(s, 10, 64); err != nil {
		return // Or: return 0, err
	}
	self.pos = idx + 1

	return
}

//fetches next string from stream and advances pos pointer
func (self *Decoder) nextString() (res string, err error) {
	if self.stream[self.pos] < '0' || self.stream[self.pos] > '9' {
		err = errors.New("No string length determinator found")
		return
	}

	//scan length
	len_start := self.pos
	len_end := self.pos
	for self.stream[len_end] != ':' {
		if len_end++; len_end >= len(self.stream) {
			err = errors.New("No string found ...")
			return
		}
	}
	len_str := string(self.stream[len_start:len_end])

	if l, e := strconv.Atoi(len_str); e != nil {
		err = fmt.Errorf("Couldn't parse string length specifier: %s", e.Error())
	} else if l >= len(self.stream[len_end:]) {
		err = errors.New("Specified length longer than data buffer ...")
	} else {
		len_end++ //skip the ':'
		res = string(self.stream[len_end : len_end+l])
		self.pos = len_end + l
	}
	return
}

//fetches a list (and its contents) from stream and advances pos
func (self *Decoder) nextList() (res []interface{}, err error) {
	if self.stream[self.pos] != 'l' {
		err = errors.New("This is not a list!")
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
func (self *Decoder) nextDict() (res map[string]interface{}, err error) {
	if self.stream[self.pos] != 'd' {
		err = errors.New("This is not a dict!")
		return
	}
	self.pos++ //skip 'd'

	res = make(map[string]interface{})

	if self.stream[self.pos] == 'e' {
		self.pos++ //skip 'e'
		return
	}

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
