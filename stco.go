package mp4box

// sample table chunk offset
/*
type stco_box struct {
	full_box_header
	count   int32    // chunk count
	entries []uint32 // chunk offset table
}
*/
func (this *encoded_box) to_stco() []uint32 {
	return this.to_uint32_slice()
}
