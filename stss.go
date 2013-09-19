package mp4box

// sample table sync sample
/*
type stss_box struct {
	full_box_header
	count   uint32
	entries []uint32 // sample index corresponds to key frame
}*/

func (this *encoded_box) to_stss() []int32 {
	return this.to_int32_slice()
}
