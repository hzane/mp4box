package mp4box

import (
	"bytes"
	"encoding/binary"
)

type smhd_box struct {
	//	full_box_header
	Balance  uint16
	Reserved uint16
}

func (this *encoded_box) to_smhd() smhd_box {
	v := smhd_box{}
	reader := bytes.NewBuffer([]byte(*this))
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	binary.Read(reader, binary.BigEndian, &v)
	return v
}
