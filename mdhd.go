package mp4box

import (
	"bytes"
	"encoding/binary"
)

// media header
// contained in mdia
// size is 24 bytes
type mdhd_v0_box struct {
	//	full_box_header  // full_box_header.Version == 0
	CreationTime     int32
	ModificationTime int32
	TimeScale        int32
	Duration         int32
	Language         uint16
	Quality          uint16
}

// size is 36 bytes
type mdhd_v1_box struct {
	//	full_box_header  // full_box_header.Version == 1
	CreationTime     int64
	ModificationTime int64
	TimeScale        int32
	Duration         int64
	Language         uint16
	Quality          uint16
}

type mdhd_box mdhd_v1_box

func (this *encoded_box) to_mdhd() mdhd_box {
	v := mdhd_box{}
	var h full_box_header
	reader := bytes.NewBuffer([]byte(*this))
	binary.Read(reader, binary.BigEndian, &h)
	switch h.Version {
	default:
		var v0 mdhd_v0_box
		binary.Read(reader, binary.BigEndian, &v0)
		v.CreationTime = int64(v0.CreationTime)
		v.ModificationTime = int64(v0.ModificationTime)
		v.TimeScale = v0.TimeScale
		v.Duration = int64(v0.Duration)
		v.Language = v0.Language
		v.Quality = v0.Quality
	case 1:
		binary.Read(reader, binary.BigEndian, &v)
	}
	return v
}
