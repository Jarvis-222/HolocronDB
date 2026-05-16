package core

import (
	"errors"
	"fmt"
	"log"
)

func readSimpleString(data []byte) (string, int, error) {

	pos := 1

	for ; data[pos] != '\r'; pos++ {

	}

	return string(data[1:pos]), pos + 2, nil

}

func readErrorMessage(data []byte) (string, int, error) {

	return readSimpleString(data)

}

func readInteger64(data []byte) (int, int, error) {

	pos := 1
	var val int64 = 0

	var sign int64 = 1

	if data[pos] == '-' {
		sign = -1
		pos += 1
	}

	for ; data[pos] != '\r'; pos++ {

		val = (val * 10) + int64(data[pos]-'0')
	}
	return int(sign * val), pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {

	len, pos, err := readInteger64(data)

	if err != nil {
		return "", 0, err
	}

	if len == -1 {
		return "", pos, nil
	}

	return string(data[pos : pos+len]), pos + len + 2, nil

}

func readArray(data []byte) ([]interface{}, int, error) {

	len, pos, err := readInteger64(data)

	if err != nil {
		return nil, 0, err
	}

	var arr []interface{} = make([]interface{}, len)
	for cnt := 0; cnt < len; cnt++ {

		val, n, err := decodeOne(data[pos:])
		arr[cnt] = val

		if err != nil {
			return nil, 0, err
		}

		pos += n
	}
	return arr, pos, nil
}

func decodeOne(data []byte) (interface{}, int, error) {

	if len(data) == 0 {
		return nil, 0, errors.New("No available data")
	}

	switch data[0] {

	case '+':
		return readSimpleString(data)

	case '-':
		return readErrorMessage(data)

	case ':':
		return readInteger64(data)

	case '$':
		return readBulkString(data)

	case '*':
		return readArray(data)

	default:
		return nil, 0, errors.New("Unsupported RESP type")
	}

}

func Decode(data []byte) (interface{}, error) {

	value, _, err := decodeOne(data)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func Encode(value interface{}, isSimpleString bool) []byte {

	log.Println("Debugging")
	switch v := value.(type) {

	case string:
		if isSimpleString {
			return fmt.Appendf(nil, "+%s\r\n", v)
		} else {
			return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(v), v)
		}

	default:
		return fmt.Appendf(nil, "%v", v)
	}

}
