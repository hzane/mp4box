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
	track_id     uint32
	duration     uint64
	volume       uint16
	width        uint32
	height       uint32
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
		log.Println(header.box_type(), header.body_size, `in mdia`)
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
		log.Println(header.box_type(), header.body_size, `in minf`)
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
	var stco, stss []uint32
	var stsz stsz_box

	foreach_child_box(stbl, func(child encoded_box, header mp4_box_header) {
		log.Println(header.box_type(), header.body_size, `in stbl`)
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
	stsz stsz_box,
	stco []uint32,
	stss []uint32) {
	log.Println(`stsz-count`, stsz.count, `stsz.size`, stsz.size)

	this.sample_count = int(stsz.count)
	this.samples = make([]mp4_sample, this.sample_count)
	for idx := 0; idx < this.sample_count; idx++ {
		if stsz.size != 0 {
			this.samples[idx].size = uint64(stsz.size)
		} else {
			this.samples[idx].size = uint64(stsz.entries[idx])
		}
	}
	this.chunk_count = len(stco)
	log.Println(this.chunk_count, `sample table chunk count`)
	this.chunks = make([]mp4_chunk, this.chunk_count)
	for idx, c := range stco {
		this.chunks[idx].offset = uint64(c)
	}

	// time to sample
	this.timestamp_count = len(stts)
	this.timestamps = make([]mp4_timestamp, this.timestamp_count)
	log.Println(this.timestamp_count, `sample table timestamp count`)

	var time_start uint64 = 0
	var start uint32 = 0
	for idx, ts := range stts {
		this.timestamps[idx].sample_start = start
		this.timestamps[idx].samples_count = ts.count
		this.timestamps[idx].time_start = time_start
		this.timestamps[idx].duration = uint64(ts.duration)
		//		log.Println(start, time_start, ts.count, `stts`, this.sample_count)
		for se := start + ts.count; start < se; start++ {
			this.samples[start].duration = uint64(ts.duration)
			this.samples[start].start_time = time_start
			time_start += uint64(ts.duration)
		}
		log.Println(`time`, time_start, `sample`, start, `timestamp to sample`)
	}

	start = 0
	for idx, sc := range stsc {
		end := this.chunk_count
		if idx < len(stsc)-1 {
			end = int(stsc[idx+1].First - 1) // first start at 1
		}
		for i := int(sc.First) - 1; i < end; i++ { // First start at 1
			this.chunks[i].sample_start = start
			this.chunks[i].sample_count = sc.SamplesPerChunk
			this.chunks[i].sample_description_id = sc.SampleDescriptionId
			var inchunk_offset uint64 = 0
			for se := start + sc.SamplesPerChunk; start < se; start++ {
				this.samples[start].chunk_id = uint32(i)
				this.samples[start].in_chunk_offset = inchunk_offset
				//				this.samples[start].description_id = sc.sample_description_id
				inchunk_offset += this.samples[start].size
			}
		}
	}

	this.sync_samples = stss
	for _, sample_id := range stss {
		this.samples[sample_id-1].is_sync_sample = true
	}
}
