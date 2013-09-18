package mp4box

// sample table sample size
type stsz_box struct {
	full_box_header
	count   uint32   // sample count
	entries []uint32 // sample size table
}

func (this *encoded_box) to_stsz() []uint32 {
	return this.to_uint32_slice()
}
