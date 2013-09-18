package mp4box

const (
	track_type_reserved = iota
	track_type_video
	track_type_audio
	track_type_hint
	track_type_other
)

type mp4_track struct {
	track_type int
}

func (this *mp4_track) from_tkhd(tkhd encoded_box) {
	tkheader := tkhd.to_tkhd()
	this.track_id = tkheader.TrackId
	this.duration = tkheader.Duration
	this.volume = tkheader.Volume
	this.width = tkheader.TrackWidth
	this.height = tkheader.TrackHeight
}

func (this *mp4_track) from_mdia(mdia encoded_box) {
	foreach_child_box(mdia, func(child encoded_box, header mp4_box_header1) {
		switch header.typ {
		case "mdhd":
			this.from_mdhd(child)
		case "minf":
			this.from_minf(child)
		case "hdlr":
			//			hdlr = child.to_hdlr()
			//			this.subtype = hdlr.subtype()
		default:
		}
	})
}

func (this *mp4_track) from_mdhd(mdhd encoded_box) {
	// do nothing
}

func (this *mp4_track) from_minf(minf encoded_box) {
	foreach_child_box(minf, func(child encoded_box, header mp4_box_header1) {
		switch header.typ {
		case "vmhd":
			this.track_type = track_type_video
		case "smhd":
			this.track_type = track_type_audio
		case "dinf":
		case "stbl":
			this.from_stbl(child)
		default:
		}
	})
}

func (this *mp4_track) from_stbl(stbl encoded_box) {
	var stsd []stsd_entry
	var stts []stts_entry
	var stsc []stsc_entry
	var stsz, stco, stss []uint32

	foreach_child_box(stbl, func(child encoded_box, header mp4_box_header1) {
		switch child.typ {
		case "stsd":
			stsd = child.to_stsd().entries
		case "stts":
			stts = child.to_stts().entries
		case "stsc":
			stsc = child.to_stsc().entries
		case "stsz":
			stsz = child.to_stsz()
		case "stco":
			stco = child.to_stco()
		case "stss":
			stss = child.to_stss()
		default:
		}
	})
	this.fill_sample_tables(stsd, stts, stsc, stsz, stco, stss)
}

func (this *mp4_track) fill_sample_tables(stsd []stsd_entry,
	stts []stts_entry,
	stsc []stsc_entry,
	stsz []uint32,
	stco []uint32,
	stss []uint32) {

	this.sample_count = len(stsz)
	this.samples = make([]mp4_sample, this.sample_count)
	for idx, ss := range stsz {
		this.samples[idx].size = ss
	}

	this.chunk_count = len(stco)
	this.chunks = make([]mp4_chunk, this.chunks_count)
	for idx, c := range stco {
		this.chunks[idx].offset = c
	}

	// time to sample
	this.timetamp_count = len(stts)
	this.timestamps = make([]timestamp_to_sample, timestamps.count)

	time_start := 0
	start = 0
	for idx, ts := range stts {
		this.timestamps[idx].sample_start = start
		this.timestamps[idx].samples_count = ts.count
		this.timestamps[idx].time_start = time_start
		this.timestamps[idx].duration = ts.duration

		for se := start + ts.sample_count; start < se; start++ {
			this.samples[start].duration = ts.duration
			this.samples[start].timestamp_start = time_start
			time_start += ts.duration
		}
	}

	/*
		First               uint32 // first chunk
		SamplesPerChunk     uint32 // samples per chunk
		SampleDescriptionId uint32 // sample description id  , index of stsd
	*/
	start = 0
	for idx, sc := range stsc {
		end := this.chunk_count
		if idx < len(stsc)-1 {
			end = stsc[idx+1].First
		}
		for i := sc.First; i < end; i++ {
			this.chunks[i].sample_begin = start
			this.chunks[i].samples_count = sc.SamplesPerChunk
			this.chunks[i].sample_description_id = sc.SampleDescriptionId
			inchunk_offset := 0
			for se := start + sc.samples_per_count; start < se; start++ {
				this.samples[start].chunk_id = i
				this.samples[start].in_chunk_offset = inchunk_offset
				//				this.samples[start].description_id = sc.sample_description_id
				inchunk_offset += this.samples[start].size
			}
		}
	}

	this.sync_samples = stss
	for _, sample_id := range stss {
		this.samples[sample_id].is_sync_sample = true
	}
}
