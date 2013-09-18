package mp4box

type vmhd_box struct {
	//	full_box_header
	GraphicsMode uint16
	Opcolor      [3]uint16
}

func (this *encoded_box) to_vmhd() vmhd_box {
	v := vmhd_box{}
	reader := bytes.NewBuffer(this)
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	binary.Read(reader, binary.BigEndian, &v)
	return v
}
