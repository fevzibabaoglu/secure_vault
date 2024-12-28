package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func EncodeDataToBytes(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeDataFromBytes(data []byte, out interface{}) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(out)
}

func EncodeInt32ToBytes(value int32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeInt32FromBytes(data []byte, out *int32) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, out)
}
