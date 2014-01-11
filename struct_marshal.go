package riaken_struct

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

import (
	"github.com/riaken/riaken-core/rpb"
)

// This package was lifted from https://github.com/mrb/riakpbc/blob/master/coder.go
// However, it was authored by myself. - Brian Jones

// typeOfBytes is a special check against Slices of []byte.
var typeOfBytes = reflect.TypeOf([]byte(nil))

// MarshalMethod is the method signature of a marshaller.
type MarshalMethod func(interface{}, *rpb.RpbContent) error

// UnmarshalMethod is the method signature of a unmarshaller.
type UnmarshalMethod func([]byte, interface{}) error

// StructMarshal contains a tag, marshaller, and unmarshaller.
// It's primary duty is to convert data from `tag` format to and from a composed struct.
type StructMarshal struct {
	tag          string          // the tag to match for the marshaller
	marshaller   MarshalMethod   // the method to run on the data
	unmarshaller UnmarshalMethod // the method to extra the data
}

// JsonMarshaller is an example of a MarshalMethod that is passed to NewEncode().
//
// If a different data marshaller is desired, such as XML, YAML, etc.  Use this as a template.
func JsonMarshaller(in interface{}, out *rpb.RpbContent) error {
	jsondata, err := json.Marshal(in)
	if err != nil {
		return err
	}
	out.Value = jsondata
	out.ContentType = []byte("application/json")
	return nil
}

// JsonUnmarshaller is an example of an UnmarshallMethod that is passed to NewEncode().
func JsonUnmarshaller(in []byte, out interface{}) error {
	err := json.Unmarshal(in, out)
	if err != nil {
		return err
	}
	return nil
}

// NewStructMarshal requires a tag and MarshalMethod.
func NewStructMarshal(tag string, marshaller MarshalMethod, unmarshaller UnmarshalMethod) *StructMarshal {
	c := new(StructMarshal)
	c.tag = tag
	c.marshaller = marshaller
	c.unmarshaller = unmarshaller
	return c
}

// Marshal takes a struct with `riak` tagged fields and builds the correct
// RpbContent to send along to Riak.
//
// Any fields of type string are set as a _bin index, and fields of any
// int type set to an _int index.
//
// Examples:
//
//  // Field is a _bin index
//  Field string `riak:"index"`
//
//  // Field is an _int index
//  Field int `riak:"index"`
//
//  // Field is a _bin index and also a json field in the actual data.
//  Field string `json:"field" riak:"index"`
func (c *StructMarshal) Marshal(data interface{}) (*rpb.RpbContent, error) {
	t := reflect.ValueOf(data)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New(fmt.Sprintf("Expected a pointer not %s", t.Kind()))
	}

	// Output
	out := &rpb.RpbContent{}

	e := t.Elem()
	switch e.Kind() {
	case reflect.Struct:
		c.process(e, out)
		if err := c.marshaller(&data, out); err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("Marshal expected a struct")
	}

	return out, nil
}

// Unmarshal unwraps the database data into the passed structure based on the defined marshaller.
func (c *StructMarshal) Unmarshal(in []byte, data interface{}) error {
	return c.unmarshaller(in, data)
}

func (c *StructMarshal) process(e reflect.Value, out *rpb.RpbContent) {
	for i := 0; i < e.NumField(); i++ {
		if !e.Field(i).CanSet() {
			continue
		}

		val := e.Field(i).Interface()
		fld := e.Type().Field(i)
		knd := e.Field(i).Kind()
		tag := fld.Tag

		// Skip anonymous fields
		if fld.Anonymous {
			continue
		}

		// Continue to process nested structs
		if knd == reflect.Struct {
			c.process(e.Field(i), out)
		}

		if tag.Get(c.tag) == "" {
			continue
		}

		if tdata := tag.Get("riak"); tdata != "" {
			for _, tfield := range strings.Split(tdata, ",") {
				switch tfield {
				case "index":
					index := &rpb.RpbPair{}
					var key string
					switch knd {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						key = fld.Name + "_int"
						switch knd {
						case reflect.Int:
							index.Value = []byte(strconv.Itoa(int(val.(int))))
							break
						case reflect.Int8:
							index.Value = []byte(strconv.Itoa(int(val.(int8))))
							break
						case reflect.Int16:
							index.Value = []byte(strconv.Itoa(int(val.(int16))))
							break
						case reflect.Int32:
							index.Value = []byte(strconv.Itoa(int(val.(int32))))
							break
						case reflect.Int64:
							index.Value = []byte(strconv.Itoa(int(val.(int64))))
							break
						}
						index.Key = []byte(strings.ToLower(key))
						break
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						key = fld.Name + "_int"
						switch knd {
						case reflect.Uint:
							index.Value = []byte(strconv.Itoa(int(val.(uint))))
							break
						case reflect.Uint8:
							index.Value = []byte(strconv.Itoa(int(val.(uint8))))
							break
						case reflect.Uint16:
							index.Value = []byte(strconv.Itoa(int(val.(uint16))))
							break
						case reflect.Uint32:
							index.Value = []byte(strconv.Itoa(int(val.(uint32))))
							break
						case reflect.Uint64:
							index.Value = []byte(strconv.Itoa(int(val.(uint64))))
							break
						}
						index.Key = []byte(strings.ToLower(key))
						break
					case reflect.String:
						key = fld.Name + "_bin"
						index.Key = []byte(strings.ToLower(key))
						index.Value = []byte(val.(string))
						break
					case reflect.Slice:
						if fld.Type == typeOfBytes {
							key = fld.Name + "_bin"
							index.Key = []byte(strings.ToLower(key))
							index.Value = []byte(val.([]byte))
						}
						break
					}
					out.Indexes = append(out.Indexes, index)
					break
				}
			}
		}
	}
}
