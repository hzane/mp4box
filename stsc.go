package mp4box

import (
	"bytes"
	"encoding/binary"
)

// sample table sample to chunk

type stsc_box struct {
	//	full_box_header
	count   uint32
	entries []stsc_entry
}

type stsc_entry struct {
	First               uint32 // first chunk index
	SamplesPerChunk     uint32 // samples per chunk
	SampleDescriptionId uint32 // sample description id  , index of stsd
}

func (this *encoded_box) to_stsc() stsc_box {
	var v stsc_box
	reader := bytes.NewBuffer([]byte(*this))
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	binary.Read(reader, binary.BigEndian, &v.count)
	v.entries = make([]stsc_entry, v.count)
	binary.Read(reader, binary.BigEndian, &v.entries)
	return v
}
