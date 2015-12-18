package lzstring

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

var testingValues map[io.Reader]string = map[io.Reader]string{
	bytes.NewBuffer([]byte{130, 33, 198, 48, 100, 2}): "abcd",
	bytes.NewBuffer([]byte{130, 55, 38, 32, 96, 9, 2, 198, 15, 225, 59, 96, 24, 2, 192, 9, 32, 158, 48, 23, 224, 11, 113, 26, 137, 163, 208, 53, 94, 5, 133, 1, 20, 80, 25, 217, 207, 110, 166, 66, 17, 196, 29, 8, 172, 205, 253, 62, 16, 167, 16, 181, 179, 1, 28, 64, 128, 214, 68, 73, 1, 216, 20, 50, 0, 139, 132, 50, 94, 0, 176, 76, 100, 6, 128, 40, 128, 43, 0, 91, 16, 66, 6, 96, 88, 112, 185, 78, 162, 253, 10, 98, 217, 80, 16, 51, 98, 230, 136, 160, 99, 130, 0, 144, 145, 221, 176, 69, 25, 0, 0, 8, 225, 29, 204, 204, 228, 16, 131, 177, 112, 113, 0, 128}): `{"dictionary":{},"dictionaryToCreate":{},"c":"","wc":"","w":"","enlargeIn":2,"dictSize":3,"numBits":2,"result":"","data":{"string":"","val":0,"position":0}}`,
}

func TestDecompress(t *testing.T) {
	for k, v := range testingValues {
		result, err := NewReaderUint16LE(k)
		if err != nil {
			t.Fatal(err)
		}

		if result.(fmt.Stringer).String() != v {
			t.Errorf(`
Expected: %s
Got:      %s`, v, result)
		}
	}
}

func BenchmarkDecompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for k := range testingValues {
			NewReaderUint16LE(k)
		}
	}
}

/*
func TestReceiveBytes(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var word uint16
		data := []uint16{}
		for e := binary.Read(r.Body, binary.LittleEndian, &word); e == nil; e = binary.Read(r.Body, binary.LittleEndian, &word) {
			data = append(data, word)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write([]byte(fmt.Sprintf("%v", data)))
	})
	if e := http.ListenAndServe(":8080", nil); e != nil {
		t.Fatal(e)
	}
}
*/
