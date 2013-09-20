package mp4box

import (
	"bytes"
	"encoding/binary"
)

// sound media header box
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

// hint media header box
type hmhd_box struct {
	// full_box_header
	max_pdu_size uint16 //gives the size in bytes of the largest PDU in this (hint) stream
	avg_pdu_size uint16 //
	max_bitrate  uint32
	avg_bitrate  uint32 //gives the maximum rate in bits/second over any window of one second
	_            uint32
}
