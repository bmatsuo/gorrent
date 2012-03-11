package main

import (
	"errors"
	"gorrent/bencode"
	"io/ioutil"
	//"bytes"
	//"fmt"
	"crypto/sha1"
)

//metainfo file (.torrent file) handling

type MetaInfo struct {
	raw    []byte
	parsed map[string]interface{}
}

func (mi *MetaInfo) ReadFromFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	mi.raw = b

	dec := bencode.NewDecoder(b)
	o, err := dec.Decode()
	if err != nil {
		return errors.New("Couldn't parse torrent: " + err.Error())
	}

	mi.parsed = o.(map[string]interface{})
	return nil
}

//return sha1 info_hash
func (mi *MetaInfo) InfoHash() []byte {
	d := mi.parsed["info"].(map[string]interface{})
	b := bencode.Encode(d)

	//sha1
	hasher := sha1.New()
	hasher.Write(b)
	//s := fmt.Sprintf("%x", hasher.Sum())
	return hasher.Sum(nil)
}
