package bencode

import (
	"fmt"
	"reflect"
	"sort"
)

//Encoder takes care of encoding objects into byte streams.
//The result of the encoding operation is available in Encoder.Bytes.
//Consecutive operations are appended to the byte stream.
//
//Accepts only string, int/int64, []interface{} and map[string]interface{} as input.
type Encoder struct {
	Bytes []byte		//the result byte stream
}

func NewEncoder() *Encoder { return new(Encoder) }

//Encode is a wrapper for Encoder.Encode.
//It returns the bencoded byte stream.
func Encode(in interface{}) []byte {
	enc := NewEncoder()
	enc.Encode(in)
	return enc.Bytes
}

//Encode encodes an object into a bencoded byte stream.
//The result of the operation is accessible through Encoder.Bytes.
//
//Example:
//	enc.Encode(23)
//	enc.Encode("test")
//	enc.Result //contains 'i23e4:test'
func (enc *Encoder) Encode(in interface{}) {
	if b := enc.encodeObject(in); len(b) > 0 {
		enc.Bytes = append(enc.Bytes, b...)
	}
}

func (enc *Encoder) encodeObject(in interface{}) []byte {
    switch t := reflect.TypeOf(in); t.Kind() {
	case reflect.String:
		return enc.encodeString(in.(string))
	case reflect.Int64:
		return enc.encodeInteger(in.(int64))
	case reflect.Int:
		return enc.encodeInteger(int64(in.(int)))
	case reflect.Slice:
		return enc.encodeList(in.([]interface{}))
	case reflect.Map:
		return enc.encodeDict(in.(map[string]interface{}))
	default:
		panic(fmt.Errorf("Can't encode this type: %s", t.Name()))
	}
	return nil
}

func (enc *Encoder) encodeString(s string) []byte {
	if len(s) <= 0 {
		return nil
	}
	return []byte(fmt.Sprintf("%d:%s", len(s), s))
}

func (enc *Encoder) encodeInteger(i int64) []byte {
	return []byte(fmt.Sprintf("i%de", i))
}

func (enc *Encoder) encodeList(list []interface{}) []byte {
	if len(list) <= 0 {
		return nil
	}
	ret := []byte("l")
    for _, obj := range list {
		ret = append(ret, enc.encodeObject(obj)...)
	}
	ret = append(ret, 'e')
	return ret
}

func (enc *Encoder) encodeDict(m map[string]interface{}) []byte {
	if len(m) <= 0 {
		return nil
	}
	//sort the map >.<
    keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ret := []byte("d")
	for _, k := range keys {
		ret = append(ret, enc.encodeString(k)...)
		ret = append(ret, enc.encodeObject(m[k])...)
	}
	ret = append(ret, 'e')
	return ret
}
