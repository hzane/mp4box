package mp4box

type mp4_sample struct {
	index           uint32
	chunk_id        uint32
	in_chunk_offset uint64 // byte
	size            uint64 // byte

	start_time     uint64 // time (time_scale)
	duration       uint64 // time_scale
	is_sync_sample bool   // sync sample / key frame
}

/*
type mp4_sample_description struct {
	typ  [4]byte
	body encoded_box
}*/
