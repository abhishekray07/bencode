package bencode

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

type encoder struct {
	bytes.Buffer
}

// Encode int64 into a byte buffer
func (m *encoder) encodeInt(number int64) {
	m.WriteByte('i')
	m.WriteString(strconv.FormatInt(number, 10))
	m.WriteByte('e')
}

// Encode uint64 into a byte buffer
func (m *encoder) encodeUint(number uint64) {
	m.WriteByte('i')
	m.WriteString(strconv.FormatUint(number, 10))
	m.WriteByte('e')
}

// Encode string into a byte buffer
func (m *encoder) encodeString(str string) {
	strLen := len(str)
	m.WriteString(strconv.Itoa(strLen))
	m.WriteByte(':')
	m.WriteString(str)
}

// encode list as a byte array
func (m *encoder) encodeList(input []interface{}) {
	m.WriteByte('l')
	for _, value := range input {
		m.chooseEncodingFunc(value)
	}
	m.WriteByte('e')
}

// encode dictionary as bytes
func (m *encoder) encodeDict(input map[string]interface{}) {
	m.WriteByte('d')
	for key, value := range input {
		// encode key as a string
		m.encodeString(key)
		m.chooseEncodingFunc(value)
	}
	m.WriteByte('e')
}

// choose the appropriate encoding function based on
// input type
func (m *encoder) chooseEncodingFunc(input interface{}) error {
	switch v := input.(type) {
	case int, int8, int16, int32, int64:
		m.encodeInt(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		m.encodeUint(reflect.ValueOf(v).Uint())
	case string:
		m.encodeString(v)
	case []interface{}:
		m.encodeList(v)
	case map[string]interface{}:
		m.encodeDict(v)
	default:
		return errors.New("Invalid type for input")
	}

	return nil
}

func encode(input interface{}) ([]byte, error) {
	m := encoder{}
	err := m.chooseEncodingFunc(input)
	if err != nil {
		return nil, err
	}

	return m.Bytes(), nil
}
