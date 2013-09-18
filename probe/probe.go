package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

var input = flag.String("input", "", "source mp4 file")

type box struct {
	typ   string
	begin int64
	size  int64
}
type strings []string

func main() {
	flag.Parse()
	if *input == "" {
		flag.PrintDefaults()
		return
	}
	f, e := os.Open(*input)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()

	var moov, mdat box

	var offset int64 = 0
	var reader io.Reader = f
	boxes := iterate_box(reader, offset, -1)
	for _, b := range boxes {
		if b.typ == "moov" {
			moov = b
		}
		if b.typ == "mdat" {
			mdat = b
		}
		fmt.Println(b)
	}
	if moov.begin <= mdat.begin {
		return
	}
	fmt.Println(moov.size, "stcos should add offset")
	boxes = iterate_stcos(reader, moov.begin, moov.begin+moov.size)
	fmt.Println(boxes)
}

func iterate_stcos(reader io.Reader, box_begin, box_end int64) []box {
	res := make([]box, 0)
	body_begin := box_begin + 8 // skip head
	boxes := iterate_box(reader, body_begin, box_end)
	var stco_names strings = []string{"stco", "co64"}
	var x strings = []string{"trak", "mdia", "minf", "stbl"}
	for _, b := range boxes {
		if x.index(b.typ) > -1 {
			res = append(res, iterate_stcos(reader, b.begin, b.begin+b.size)...)
		} else if stco_names.index(b.typ) > -1 {
			res = append(res, b)
		}
	}
	return res
}

func (strs strings) index(o string) int {
	for idx, str := range strs {
		if str == o {
			return idx
		}
	}
	return -1
}
func iterate_box(reader io.Reader, begin, end int64) (boxes []box) {
	boxes = make([]box, 0)
	//	seeker := reader.(io.Seeker)
	//	seeker.Seek(begin, 0)

	for box, err := read_next_box_header(reader, begin); err == nil && (end < 0 || begin < end); box, err = read_next_box_header(reader, begin) {
		if box.size == 0 {
			break
		}
		boxes = append(boxes, box)
		begin += box.size
	}
	return
}

func read_next_box_header(reader io.Reader, begin int64) (b box, e error) {
	seeker := reader.(io.Seeker)
	b.begin, e = seeker.Seek(begin, 0)

	var sz int32
	t := make([]byte, 4)
	e = binary.Read(reader, binary.BigEndian, &sz)
	_, e = reader.Read(t)
	b.typ = string(t)
	if sz == 0 { //end box of mdat
		return
	}
	if sz == 1 {
		var sz64 int64
		e = binary.Read(reader, binary.BigEndian, &sz64)
		b.size = sz64
		return
	}
	b.size = int64(sz)
	return
}
