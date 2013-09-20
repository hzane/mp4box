package mp4box

import (
	"bytes"
	"encoding/binary"
)

// sample table sample size
type stsz_box struct {
	//	full_box_header
	size    uint32   // sample size if all sample has the same size
	count   uint32   // sample count
	entries []uint32 // sample size table
}

func (this *encoded_box) to_stsz() stsz_box {
	reader := bytes.NewBuffer([]byte(*this))
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	v := stsz_box{}
	binary.Read(reader, binary.BigEndian, &v.size)
	binary.Read(reader, binary.BigEndian, &v.count)
	v.entries = make([]uint32, v.count)
	binary.Read(reader, binary.BigEndian, &v.entries)
	return v
}
