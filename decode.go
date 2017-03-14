package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type decoder struct {
	bufio.Reader
}

// read message until the specified byte
// returns the message as a string
func (m *decoder) readUntil(until byte) (interface{}, error) {
	res, err := m.ReadSlice(until)
	if err != nil {
		return nil, err
	}

	resultStr := string(res[:len(res)-1])
	return resultStr, nil
}

// selects the appropriate decode function based on the
// first byte of the message
func (m *decoder) chooseDecodeFunc(first []byte) (interface{}, error) {
	switch string(first) {
	case "i":
		m.ReadByte()
		return m.decodeInt()
	case "d":
		m.ReadByte()
		return m.decodeDict()
	case "e":
		return nil, errors.New("Invalid character e. This should be handled by calling function")
	default:
		return m.decodeString()
	}
}

// checks if the first byte is as expected
func (m *decoder) checkByte(expected byte) bool {
	firstByte, err := m.ReadByte()
	if err != nil {
		return false
	}

	if firstByte != expected {
		return false
	}

	return true
}

// decode message as a string
// returns the string as a byte array
func (m *decoder) decodeString() (interface{}, error) {
	res, err := m.readUntil(':')
	if err != nil {
		return nil, err
	}

	var data string
	var ok bool

	if data, ok = res.(string); !ok {
		return nil, errors.New("Error during decodeString. Incorrect Format")
	}

	// get length of the string
	strLen, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, err
	}

	// read the string from the buffer
	byteStr := make([]byte, strLen)
	_, err = io.ReadFull(m, byteStr)
	if err != nil {
		return nil, err
	}

	return byteStr, nil
}

// decode message as int
func (m *decoder) decodeInt() (interface{}, error) {
	res, err := m.readUntil('e')
	if err != nil {
		return nil, err
	}

	var data string
	var ok bool

	if data, ok = res.(string); !ok {
		return nil, errors.New("Error during decodeInt. Incorrect format")
	}

	// get length of the string
	num, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, err
	}

	return num, nil
}

// decode message as list
func (m *decoder) decodeList() (interface{}, error) {
	checkFirst := m.checkByte('l')
	if !checkFirst {
		return nil, errors.New("Error during decodeList")
	}

	var list []interface{}

	for {
		nextByte, err := m.Peek(1)
		if err != nil {
			return "", nil
		}

		if string(nextByte) == "e" {
			m.ReadByte()
			break
		} else {
			res, err := m.chooseDecodeFunc(nextByte)
			if err != nil {
				return nil, err
			}

			list = append(list, res)
		}
	}
	return list, nil
}

// decode the message as a dictionary.
// the key of the dictionary is a string
func (m *decoder) decodeDict() (interface{}, error) {
	checkFirst := m.checkByte('d')
	if !checkFirst {
		return nil, errors.New("Error during decodeDict")
	}

	dict := make(map[string]interface{})

	for {
		nextByte, err := m.Peek(1)
		if err != nil {
			return "", nil
		}

		if string(nextByte) == "e" {
			m.ReadByte()
			break
		} else {
			res, err := m.chooseDecodeFunc(nextByte)
			if err != nil {
				return nil, err
			}

			var byteKey []byte
			var ok bool

			if byteKey, ok = res.([]byte); !ok {
				return nil, errors.New("Dictionary key is not a valid string")
			}

			key := string(byteKey)
			val, err := m.chooseDecodeFunc(nextByte)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}
	}

	return dict, nil
}

// decode bencoded string
func decode(reader io.Reader) (interface{}, error) {
	m := decoder{*bufio.NewReader(reader)}
	firstByte, err := m.Peek(1)
	if err != nil {
		return "", err
	}

	if string(firstByte) != "d" {
		return nil, errors.New("Bencoded string doesn't start with d")
	}

	return m.decodeDict()
}
