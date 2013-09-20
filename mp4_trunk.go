package mp4box

type mp4_chunk struct {
	sample_start          uint32
	sample_count          uint32
	sample_description_id uint32

	time_start     uint64 // timescale
	duration       uint64 // timescale
	offset         uint64 // file position
	length         uint64 // bytes
	offset_in_mdat uint64
}
