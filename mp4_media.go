package mp4box

type mp4_media struct {
	//	mp4_header
	video_track mp4_track // has desc
	audio_track mp4_track // has desc

	video_desc mp4_sample_description
	audio_desc mp4_sample_description

	samples      []mp4_sample
	chunks       []mp4_chunk
	timestamps   []mp4_timestamp
	sync_samples []uint32

	mdat_offset uint64 // in file position byte
	ftyp_offset uint64 // in file
	moov_offset uint64 // in line
	brand       string
}

func new_mp4_media(reader io.ReadSeeker) (fd mp4_media, err error) {
	h := next_box_header(reader)
	ftyp := next_box_body(reader, h).to_ftyp()
	fd.brand = string(ftyp.major_brand)

	for {
		h = next_box_header(reader)
		if h.size == 0 {
			break
		}
		switch h.typ {
		default:
			reader.Seek(h.body_size, 1)
		case "moov":
			fd.moov_offset = reader.Seek(0, 0) - (h.size - h.body_size)
			fd.moov_body_length = h.body_size()
			fd.from_moov(next_box_body(reader, h))
			break
		case "mdat":
			fd.mdat_offset = reader.Seek(0, 0) - (h.size - h.body_size)
			fd.mdat_body_length = h.body_size()
			reader.Seek(h.body_size, 1)
		}
	}
	return
}

func (this *mp4_media) from_mvhd(mvhd encoded_box) {
	mvheader := mvhd.to_mvhd()
	this.time_scale = mvheader.TimeScale
	this.duration = mvheader.Duration
	this.rate = mvheader.Rate
	this.volume = mvheader.Volume
}

func (this *mp4_media) from_moov(moov encoded_box) {
	foreach_child_box(moov, func(child encoded_box, header mp4_box_header) {
		switch header.typ {
		default:
		case "mvhd":
			fd.from_mvhd(child)
		case "trak":
			this.from_trak(child)
		}
	})
}

func (this *mp4_media) from_trak(trak encoded_box) {
	var t mp4_track
	foreach_child_box(trak, func(child encoded_box, header mp4_box_header) {
		switch header.typ {
		case "tkhd":
			t.from_tkhd(child)
		case "mdia":
			t.from_mdia(child)
		default:
		}
	})
	switch t.track_type {
	default:
	case track_type_audio:
		this.audio_track = t
	case track_type_video:
		this.video_track = t
	case track_type_hint:
	}
}
