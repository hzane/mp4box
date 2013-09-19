package mp4box

import (
	"log"
)

const (
	track_type_reserved = iota
	track_type_video
	track_type_audio
	track_type_hint
	track_type_other
)

type mp4_track struct {
	track_type   int
	track_id     int32
	duration     int64
	volume       uint16
	width        int32
	height       int32
	sample_count int
	samples      []mp4_sample

	chunk_count int
	chunks      []mp4_chunk

	timestamp_count int
	timestamps      []mp4_timestamp

	sync_samples []uint32
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
	foreach_child_box(mdia, func(child encoded_box, header mp4_box_header) {
		log.Println(header.typ, header.body_size, `------`)
		switch header.box_type() {
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
	foreach_child_box(minf, func(child encoded_box, header mp4_box_header) {
		log.Println(header.typ, header.body_size, `			`)
		switch header.box_type() {
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

	foreach_child_box(stbl, func(child encoded_box, header mp4_box_header) {
		log.Println(header.typ, header.body_size, `			`)
		switch header.box_type() {
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
		this.samples[idx].size = int64(ss)
	}

	this.chunk_count = len(stco)
	this.chunks = make([]mp4_chunk, this.chunk_count)
	for idx, c := range stco {
		this.chunks[idx].offset = int64(c)
	}

	// time to sample
	this.timestamp_count = len(stts)
	this.timestamps = make([]mp4_timestamp, this.timestamp_count)

	var time_start int64 = 0
	var start int32 = 0
	for idx, ts := range stts {
		this.timestamps[idx].sample_start = start
		this.timestamps[idx].samples_count = ts.Count
		this.timestamps[idx].time_start = time_start
		this.timestamps[idx].duration = int64(ts.Duration)

		for se := start + ts.Count; start < se; start++ {
			this.samples[start].duration = int64(ts.Duration)
			this.samples[start].start_time = time_start
			time_start += int64(ts.Duration)
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
			end = int(stsc[idx+1].First)
		}
		for i := int(sc.First); i < end; i++ {
			this.chunks[i].sample_start = start
			this.chunks[i].sample_count = sc.SamplesPerChunk
			this.chunks[i].sample_description_id = sc.SampleDescriptionId
			var inchunk_offset int64 = 0
			for se := start + sc.SamplesPerChunk; start < se; start++ {
				this.samples[start].chunk_id = int32(i)
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
