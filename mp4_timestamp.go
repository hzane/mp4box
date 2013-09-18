package mp4box

type mp4_timestamp struct {
	sample_start  uint32
	samples_count uint32
	time_start    uint64
	duration      uint64
}
