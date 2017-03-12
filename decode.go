package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

/*

func readUntil

func decodeString

func decodeInt

func decodeList

func decodeDictionary

func decode

*/

type message struct {
	bufio.Reader
}

// read message until the specified byte
// returns the message as a string
func (m *message) readUntil(until byte) (interface{}, error) {
	res, err := m.ReadSlice(until)
	if err != nil {
		return nil, err
	}

	resultStr := string(res[:len(res)-1])
	return resultStr, nil
}

// decode message as a string
func (m *message) decodeString() (interface{}, error) {
	res, err := m.readUntil(':')
	if err != nil {
		return nil, err
	}

	var data string
	var ok bool

	if data, ok = res.(string); !ok {
		return nil, errors.New("Error during decodeString. Incorrect Format")
	}

	strLen, err := strconv.Atoi(data)
	if err != nil {
		return nil, err
	}

	str := make([]byte, strLen)
	_, err = io.ReadFull(m, str)
	if err != nil {
		return nil, err
	}

	return string(str), nil
}

// decode message as int
func (m *message) decodeInt() (interface{}, error) {
	res, err := m.readUntil('e')
	if err != nil {
		return nil, err
	}

	var str string
	var ok bool

	if str, ok = res.(string); !ok {
		return nil, errors.New("Error during decodeInt. Incorrect format")
	}

	return str[1:], nil
}

// decode message as list
func (m *message) decodeList() (interface{}, error) {
	var list []interface{}

	firstByte, err := m.ReadByte()
	if err != nil {
		return nil, err
	}

	if firstByte != 'l' {
		return nil, errors.New("Error decoding list. First character is not l")
	}

	for {
		nextByte, err := m.Peek(1)
		if err != nil {
			return "", nil
		}

		switch string(nextByte) {
		case "e":
			m.ReadByte()
			return list, nil
		case "i":
			res, _ := m.decodeInt()
			list = append(list, res)
			fmt.Println("Decode int 2", res)
		default:
			res, _ := m.decodeString()
			list = append(list, res)
			fmt.Println("Decode String 2", res)
		}
	}
}

func (m *message) decodeDict() (interface{}, error) {
	dict := make(map[interface{}]interface{})
	var key interface{}

	firstByte, err := m.ReadByte()
	if err != nil {
		return nil, err
	}

	if firstByte != 'd' {
		return nil, errors.New("Error decoding dictionary. First character is not d")
	}

	key = nil

	for {
		fmt.Println("Key = ", key)

		nextByte, err := m.Peek(1)
		if err != nil {
			return "", nil
		}

		switch string(nextByte) {
		case "e":
			m.ReadByte()
			return dict, nil
		case "l":
			res, _ := m.decodeList()
			dict[key] = res
		case "i":
			res, _ := m.decodeInt()
			dict[key] = res
		default:
			res, _ := m.decodeString()
			fmt.Println("res", res)
			if key == nil {
				key = res
			} else {
				dict[key] = res
				key = nil
			}
		}
	}
}

func decode(reader io.Reader) (interface{}, error) {
	m := message{*bufio.NewReader(reader)}

	for {
		firstByte, err := m.Peek(1)
		if err != nil {
			return "", nil
		}

		switch string(firstByte) {
		case "i":
			res, _ := m.decodeInt()
			fmt.Println("Decode int", res)
		case "d":
			res, _ := m.decodeDict()
			fmt.Println("Decode dict", res)
		case "l":
			res, _ := m.decodeList()
			fmt.Println("Decode list", res)
		default:
			res, _ := m.decodeString()
			fmt.Println("Decode String", res)
		}
	}
}

func main() {
	str := "d3:cow3:moo4:spam4:eggse"
	buf := bytes.NewBufferString(str)

	decode(buf)
}
