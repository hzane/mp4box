package mp4box

type tkhd_v0_box struct {
	// full_box_header
	CreationTime     uint32
	ModificationTime uint32
	TrackId          uint32
	Reserved         uint32
	Duration         uint32
	Reserved1        uint64
	Layer            uint16
	AltermateGrouop  uint16
	Volume           uint16
	Reserved2        uint16
	Matrix           [36]byte
	TrackWidth       uint32
	TrackHeight      uint32
}

type tkhd_v1_box struct {
	// full_box_header
	CreationTime     uint64
	ModificationTime uint64
	TrackId          uint32 // not zero
	Reserved         uint32
	Duration         uint64 // timescale
	Reserved1        uint64
	Layer            uint16
	AltermateGrouop  uint16
	Volume           uint16 // 8.8 fixed-float 1.0 means normal
	Reserved2        uint16
	Matrix           [36]byte
	TrackWidth       uint32
	TrackHeight      uint32
}
type tkhd_box tkhd_v1_box

func (this *encoded_box) to_tkhd() tkhd_box {
	v := tkhd_box{}

	reader := bytes.NewBuffer(this)
	var h full_box_header
	binary.Read(reader, binary.BigEndian, &h)
	switch h.Version {
	default:
		var v0 tkhd_v0_box
		binary.Read(reader, binary.BigEndian, &v0)
		v.CreationTime = uint64(v0.CreationTime)
		v.ModificationTime = uint64(v0.ModificationTime)
		v.Track_id = v0.Track_id
		v.Reserved = v0.Reserved
		v.Duration = uint64(v0.Duration)
		v.Reserved1 = v0.Reserved1
		v.Layer = v0.Layer
		v.AltermateGrouop = v0.AltermateGrouop
		v.Volume = v0.Volume
		v.Reserved2 = v0.Reserved2
		copy(v.Matrix, v0.Matrix)
		v.TrackWidth = v0.TrackWidth
		v.TrackHeight = v0.TrackHeight
	case 1:
		binary.Read(reader, binary.BigEndian, &v)
	}
	return v
}
