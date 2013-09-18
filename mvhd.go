package mp4box

type mvhd_v0_box struct {
	// full_box_header
	CreationTime     uint32 // unix time
	ModificationTime uint32 // unix time
	TimeScale        uint32
	Duration         uint32
	Rate             uint32 //[16.16] 0x00010000 = 1.0
	Volume           uint16 //[8.8] 0x0100 = 1.0
	Reserved         [10]byte
	Matrix           [36]byte
	Predefined       [24]byte
	NextTrack        uint32
}

type mvhd_v1_box struct {
	// full_box_header
	CreationTime     uint64 // unix time
	ModificationTime uint64 // unix time
	TimeScale        uint32
	Duration         uint64
	Rate             uint32 //[16.16] 0x00010000 = 1.0
	Volume           uint16 //[8.8] 0x0100 = 1.0
	Reserved         [10]byte
	Matrix           [36]byte
	Predefined       [24]byte
	NextTrack        uint32
}
type mvhd_box mvhd_v1_box

func (this *encoded_box) to_mvhd() mvhd_box {
	v := mvhd_box{}
	buf := bytes.NewBuffer(this)
	var h full_box_header
	binary.Read(buf, binary.BigEndian, &h)
	switch h.Version {
	default:
		var v0 mvhd_v0_box
		binary.Read(buf, binary.BigEndian, &v0)
		v.CreationTime = uint64(v0.CreationTime)
		v.ModificationTime = uint64(v0.ModificationTime)
		v.TimeScale = v0.TimeScale
		v.Duration = uint64(v0.Duration)
		v.Rate = v0.Rate
		v.Volume = v0.Volume
		copy(v.Reserved, v0.Reserved)
		copy(v.Matrix, v0.Matrix)
		copy(v.Predefined, v0.Predefined)
		v.NextTrack = v0.NextTrack
	case 1:
		binary.Read(buf, binary.BigEndian, &v)
	}
	return v
}
