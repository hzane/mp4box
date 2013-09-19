package mp4box

type mp4_sample struct {
	index           int32
	chunk_id        int32
	in_chunk_offset int64 // byte
	size            int64 // byte

	start_time     int64 // time (time_scale)
	duration       int64 // time_scale
	is_sync_sample bool  // sync sample / key frame
}

type mp4_sample_description struct {
	typ  [4]byte
	body encoded_box
}
