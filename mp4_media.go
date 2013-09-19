package mp4box

import (
	"io"
	"log"
)

type Mp4Media struct {
	//	mp4_header
	video_track mp4_track // has desc
	audio_track mp4_track // has desc

	video_desc mp4_sample_description
	audio_desc mp4_sample_description

	samples      []mp4_sample
	chunks       []mp4_chunk
	timestamps   []mp4_timestamp
	sync_samples []uint32

	mdat_offset      int64 // in file position byte
	ftyp_offset      int64 // in file
	moov_offset      int64 // in line
	moov_body_length int64
	mdat_body_length int64
	time_scale       int64
	duration         int64
	volume           uint16
	rate             int32

	brand string
}

func NewMp4Media(reader io.ReadSeeker) (fd *Mp4Media, err error) {
	fd = &Mp4Media{}
	h := next_box_header(reader)
	ftyp := next_box_body(reader, h).to_ftyp()
	fd.brand = string(ftyp.major_brand[:])

	for {
		h = next_box_header(reader)
		if h.size == 0 {
			break
		}
		log.Println(h.typ, h.body_size)
		switch string(h.typ[:]) {
		default:
			reader.Seek(h.body_size, 1)
		case "moov":
			mo, _ := reader.Seek(0, 0)
			fd.moov_offset = mo - (h.size - h.body_size)
			fd.moov_body_length = h.body_size
			fd.from_moov(next_box_body(reader, h))
			break
		case "mdat":
			mo, _ := reader.Seek(0, 0)
			fd.mdat_offset = mo - (h.size - h.body_size)
			fd.mdat_body_length = h.body_size
			reader.Seek(h.body_size, 1)
		}
	}
	return
}

func (this *Mp4Media) from_mvhd(mvhd encoded_box) {
	mvheader := mvhd.to_mvhd()
	this.time_scale = int64(mvheader.TimeScale)
	this.duration = int64(mvheader.Duration)
	this.rate = mvheader.Rate
	this.volume = mvheader.Volume
}

func (this *Mp4Media) from_moov(moov encoded_box) {
	foreach_child_box(moov, func(child encoded_box, header mp4_box_header) {
		log.Println(header.typ, header.body_size, `	`)
		switch header.box_type() {
		default:
		case "mvhd":
			this.from_mvhd(child)
		case "trak":
			this.from_trak(child)
		}
	})
}

func (this *Mp4Media) from_trak(trak encoded_box) {
	var t mp4_track
	foreach_child_box(trak, func(child encoded_box, header mp4_box_header) {
		log.Println(header.typ, header.size, `		`)
		switch header.box_type() {
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
