package bencode

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
)

// Value может быть int, string, list или dict
type Value interface {
	isValue()
}

type (
	Int    int64
	String []byte
	List   []Value
	Dict   map[string]Value
)

func (Int) isValue()    {}
func (String) isValue() {}
func (List) isValue()   {}
func (Dict) isValue()   {}

// Unmarshal десериализует bencode данные
func Unmarshal(data []byte) (Value, error) {
	reader := bytes.NewReader(data)
	return decodeValue(reader)
}

// UnmarshalTo маппит bencode в структуру
func UnmarshalTo(data []byte, v interface{}) error {
	val, err := Unmarshal(data)
	if err != nil {
		return err
	}
	return mapValue(val, v)
}

// Marshal сериализует в bencode
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := encodeValue(&buf, v)
	return buf.Bytes(), err
}

// decodeValue читает одно значение из потока
func decodeValue(r *bytes.Reader) (Value, error) {
	typ, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {
	case typ == 'i': // integer
		return decodeInt(r)
	case typ == 'l': // list
		return decodeList(r)
	case typ == 'd': // dict
		return decodeDict(r)
	case typ >= '0' && typ <= '9': // string
		return decodeStringWithLength(r, int(typ-'0'))
	default:
		return nil, fmt.Errorf("unknown bencode type: %c", typ)
	}
}

// decodeStringWithLength читает строку, когда длина уже прочитана как первый символ
func decodeStringWithLength(r *bytes.Reader, firstDigit int) (Value, error) {
	// Читаем остальную часть длины (до ':')
	var lenBuf bytes.Buffer
	lenBuf.WriteByte(byte('0' + firstDigit))
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == ':' {
			break
		}
		lenBuf.WriteByte(b)
	}

	length, err := strconv.ParseInt(lenBuf.String(), 10, 64)
	if err != nil {
		return nil, err
	}

	// Читаем саму строку
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	return String(data), nil
}

func decodeInt(r *bytes.Reader) (Value, error) {
	var buf bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 'e' {
			break
		}
		buf.WriteByte(b)
	}

	str := buf.String()
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, err
	}
	return Int(val), nil
}

func decodeString(r *bytes.Reader) (Value, error) {
	// Читаем длину (до ':')
	var lenBuf bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == ':' {
			break
		}
		lenBuf.WriteByte(b)
	}

	length, err := strconv.ParseInt(lenBuf.String(), 10, 64)
	if err != nil {
		return nil, err
	}

	// Читаем саму строку
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	return String(data), nil
}

func decodeList(r *bytes.Reader) (Value, error) {
	list := List{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 'e' {
			break
		}
		r.UnreadByte()

		val, err := decodeValue(r)
		if err != nil {
			return nil, err
		}
		list = append(list, val)
	}
	return list, nil
}

func decodeDict(r *bytes.Reader) (Value, error) {
	dict := Dict{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 'e' {
			break
		}
		r.UnreadByte()

		// Ключ должен быть строкой
		keyVal, err := decodeValue(r)
		if err != nil {
			return nil, fmt.Errorf("decodeDict: failed to decode key: %w", err)
		}
		key, ok := keyVal.(String)
		if !ok {
			return nil, errors.New("dict key must be string")
		}

		// Значение (не делаем UnreadByte здесь — decodeValue читает сам)
		val, err := decodeValue(r)
		if err != nil {
			return nil, fmt.Errorf("decodeDict: failed to decode value for key %s: %w", key, err)
		}

		dict[string(key)] = val
	}
	return dict, nil
}

// encodeValue сериализует значение в bencode
func encodeValue(w io.Writer, v interface{}) error {
	switch val := v.(type) {
	case Int:
		fmt.Fprintf(w, "i%de", val)
	case String:
		fmt.Fprintf(w, "%d:%s", len(val), val)
	case List:
		if _, err := w.Write([]byte("l")); err != nil {
			return err
		}
		for _, item := range val {
			if err := encodeValue(w, item); err != nil {
				return err
			}
		}
		_, err := w.Write([]byte("e"))
		return err
	case Dict:
		if _, err := w.Write([]byte("d")); err != nil {
			return err
		}
		// Сортируем ключи для канонической формы
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if err := encodeValue(w, String(k)); err != nil {
				return err
			}
			if err := encodeValue(w, val[k]); err != nil {
				return err
			}
		}
		_, err := w.Write([]byte("e"))
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

// mapValue маппит bencode Value в Go структуру
func mapValue(val Value, dest interface{}) error {
	// TODO: реализовать маппинг в структуры
	return nil
}

// GetString получает строку из Value
func GetString(v Value, key string) (string, bool) {
	if dict, ok := v.(Dict); ok {
		if val, exists := dict[key]; exists {
			if str, ok := val.(String); ok {
				return string(str), true
			}
		}
	}
	return "", false
}

// GetInt получает int из Value
func GetInt(v Value, key string) (int64, bool) {
	if dict, ok := v.(Dict); ok {
		if val, exists := dict[key]; exists {
			if i, ok := val.(Int); ok {
				return int64(i), true
			}
		}
	}
	return 0, false
}

// GetDict получает dict из Value
func GetDict(v Value, key string) (Dict, bool) {
	if dict, ok := v.(Dict); ok {
		if val, exists := dict[key]; exists {
			if d, ok := val.(Dict); ok {
				return d, true
			}
		}
	}
	return nil, false
}

// GetList получает list из Value
func GetList(v Value, key string) (List, bool) {
	if dict, ok := v.(Dict); ok {
		if val, exists := dict[key]; exists {
			if l, ok := val.(List); ok {
				return l, true
			}
		}
	}
	return nil, false
}

// ToDict конвертирует в Dict если возможно
func ToDict(v Value) (Dict, bool) {
	if d, ok := v.(Dict); ok {
		return d, true
	}
	return nil, false
}

// Bytes возвращает строку как []byte
func (s String) Bytes() []byte {
	return []byte(s)
}

// String возвращает значение как строку
func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Binary helpers для работы с большими числами
func Uint64ToInt(u uint64) Int {
	return Int(u)
}

func IntToUint64(i Int) uint64 {
	return uint64(i)
}

// ReadUint64 читает uint64 из reader (big-endian)
func ReadUint64(r io.Reader) (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf[:]), nil
}

// WriteUint64 записывает uint64 в writer (big-endian)
func WriteUint64(w io.Writer, v uint64) error {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}

// DecodeHexString декодирует hex строку в байты
func DecodeHexString(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// EncodeHexString кодирует байты в hex строку
func EncodeHexString(data []byte) string {
	return hex.EncodeToString(data)
}
