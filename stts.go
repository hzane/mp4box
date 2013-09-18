package mp4box

// sample table time to sample map
type stts_box struct {
	//	full_box_header
	count   uint32
	entries []stts_entry
}

// timestamp to sample
type stts_entry struct {
	Count    uint32 // sample count
	Duration uint32 // sample duration
}

func (this *encoded_box) to_stts() stts_box {
	reader := bytes.NewBuffer(this)
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	var v stts_box
	binary.Read(reader, binary.BigEndian, &v.count)
	v.entries = make([]stts_entry, v.count)
	for i := 0; i < v.count; i++ {
		binary.Read(reader, binary.BigEndian, &v.entries[i])
	}
	return v
}
