package lzstring

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
)

func NewReader(r io.Reader) (io.ReadCloser, error) {
	d := &reader{
		position:  32768,
		data:      bufio.NewReader(r),
		dict:      map[int]string{0: "0", 1: "1", 2: "2"},
		dictSize:  4,
		result:    &bytes.Buffer{},
		numBits:   3,
		enlargeIn: 4,
	}
	if e := d.process(); e != nil {
		return nil, e
	}
	return d, nil
}

type reader struct {
	currentVal                               int32
	position                                 int32
	data                                     *bufio.Reader
	dict                                     map[int]string
	lastEntry                                string
	result                                   *bytes.Buffer
	numBits, errorCount, dictSize, enlargeIn int
}

func (r *reader) readBit() int32 {
	res := r.currentVal & r.position
	r.position = r.position >> 1
	if r.position == 0 {
		r.position = 32768

		if val, _, e := r.data.ReadRune(); e == io.EOF {
			r.currentVal = 0
		} else if e != nil {
			// TODO: Deal with error
		} else {
			r.currentVal = val
		}
	}
	if res > 0 {
		return 1
	}
	return 0
}

func (r *reader) readBits(numBits int) (res int32) {
	maxpower := math.Pow(2, float64(numBits))
	power := 1
	for float64(power) != maxpower {
		res = res | (int32(power) * r.readBit())
		power = power << 1
	}
	return
}

func (r *reader) Close() error {
	return nil
}

func (r *reader) process() error {
	if v, _, e := r.data.ReadRune(); e != nil {
		return e
	} else {
		r.currentVal = v
	}

	var c int32
	switch r.readBits(2) {
	case 0:
		c = r.readBits(8)
	case 1:
		c = r.readBits(16)
	case 2:
		return nil
	}

	r.dict[3] = fmt.Sprintf("%c", c)
	r.result.WriteRune(c)
	r.lastEntry = r.dict[3]

	for {

		bitsRead := r.readBits(r.numBits)
		var 			dictIndex = r.dictSize
		
		switch bitsRead {
		case 0:
			if r.errorCount++; r.errorCount > 10000 {
				return fmt.Errorf("Error count exceded 10000")
			}
			c = r.readBits(8)
			r.dict[r.dictSize] = fmt.Sprintf("%c", c)
			r.dictSize++
			r.enlargeIn--

		case 1:
			c = r.readBits(16)
			r.dict[r.dictSize] = fmt.Sprintf("%c", c)
			r.dictSize++
			r.enlargeIn--

		case 2:
			return nil

		default:
			dictIndex = int(bitsRead)
		}

		if r.enlargeIn == 0 {
			r.enlargeIn = int(math.Pow(2, float64(r.numBits)))
			r.numBits++
		}

		v, exists := r.dict[dictIndex]
		if !exists {
			if dictIndex == r.dictSize {
				v = fmt.Sprintf("%s%c", r.lastEntry, []rune(r.lastEntry)[0])
			} else {
				return nil
			}
		}
		r.result.WriteString(v)
		r.dict[r.dictSize] = fmt.Sprintf("%s%c", r.lastEntry, []rune(v)[0])
		r.dictSize++
		r.lastEntry = v

		if r.enlargeIn--; r.enlargeIn == 0 {
			r.enlargeIn = int(math.Pow(2, float64(r.numBits)))
			r.numBits++
		}
	}

	return nil
}

func (r *reader) Read(b []byte) (int, error) {
	return r.result.Read(b)
}

func (r *reader) String() string {
	return r.result.String()
}
