package riaken_struct

import (
	"encoding/json"
	"testing"
)

type EncodeData struct {
	Email        string     `json:"email" riak:"index"`
	Twitter      string     `json:"twitter" riak:"index"`
	Data         []byte     `json:"data" riak:"index"`
	AnInt        int        `json:"anint" riak:"index"`
	AnInt8       int8       `json:"anint8" riak:"index"`
	AnInt16      int16      `json:"anint16" riak:"index"`
	AnInt32      int32      `json:"anint32" riak:"index"`
	AnInt64      int64      `json:"anint64" riak:"index"`
	AUInt        uint       `json:"auint" riak:"index"`
	AUInt8       uint8      `json:"auint8" riak:"index"`
	AUInt16      uint16     `json:"auint16" riak:"index"`
	AUInt32      uint32     `json:"auint32" riak:"index"`
	AUInt64      uint64     `json:"auint64" riak:"index"`
	Byte         byte       `json:"abyte" riak:"index"`
	Rune         rune       `json:"arune" riak:"index"`
	NestedStruct NestedData `json:"nested_struct"`
	unexported   string
}

type NestedData struct {
	Nested string `json:"nested" riak:"index"`
}

func TestCoder(t *testing.T) {
	e := NewStructMarshal("json", JsonMarshaller, JsonUnmarshaller)

	data := &EncodeData{
		Email:   "riak@example.com",
		Twitter: "riak-twitter",
		Data:    []byte("riak-data"),
		AnInt:   1,
		AnInt8:  127,
		AnInt16: 32767,
		AnInt32: 2147483647,
		AnInt64: 9223372036854775807,
		AUInt:   1,
		AUInt8:  255,
		AUInt16: 65535,
		AUInt32: 4294967295,
		AUInt64: 18446744073709551615,
		Byte:    255,
		Rune:    2147483647,
		NestedStruct: NestedData{
			Nested: "nested-data",
		},
		unexported: "unexported",
	}

	encdata, err := e.Marshal(data)
	if err != nil {
		t.Error(err.Error())
	}

	key := string(encdata.GetIndexes()[0].GetKey())
	if key != "email_bin" {
		t.Errorf("Expected email_bin, got %s", key)
	}

	var found bool
	for _, v := range encdata.GetIndexes() {
		if string(v.GetKey()) == "nestedstruct_nested_bin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected a nestedstruct_nested_bin index")
	}

	jsondata, err := json.Marshal(data)
	if err != nil {
		t.Error(err.Error())
	}
	if string(encdata.GetValue()) != string(jsondata) {
		t.Errorf("Expected %s, got %s", string(encdata.GetValue()), string(jsondata))
	}

	result := &EncodeData{}
	e.Unmarshal(jsondata, result)
	if result.Email != data.Email {
		t.Errorf("Expected %s, got %s", data.Email, result.Email)
	}
}
