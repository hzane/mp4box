package mp4box

type mp4_chunk struct {
	sample_start          int32
	sample_count          int32
	sample_description_id int32

	time_start     int64 // timescale
	duration       int64 // timescale
	offset         int64 // file position
	length         int64 // bytes
	offset_in_mdat int64
}
