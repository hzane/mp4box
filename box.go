package mp4box

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

type encoded_box []byte
type mp4_box_header struct {
	size      int64
	typ       [4]byte
	body_size int64
}

func (this mp4_box_header) box_type() string {
	return string(this.typ[:])
}
func next_box_header(reader io.Reader) (v mp4_box_header) {
	var mp4_box_header_default_size int64 = 8 // sizeof(size) + sizeof(type)
	var size uint32
	binary.Read(reader, binary.BigEndian, &size)
	switch size {
	default:
		v.size = int64(size)
		v.body_size = v.size - mp4_box_header_default_size
		binary.Read(reader, binary.BigEndian, &v.typ)
	case 1:
		binary.Read(reader, binary.BigEndian, &v.typ)
		binary.Read(reader, binary.BigEndian, &v.size)
		v.body_size = v.size - mp4_box_header_default_size - 8 // 8 == sizeof(Extend)
	case 0:
	}

	return
}

func next_box_body(reader io.Reader, h mp4_box_header) encoded_box {
	x := make(encoded_box, h.body_size)
	reader.Read(x)
	return x
}

func foreach_child_box(b encoded_box, f func(child encoded_box, h mp4_box_header)) {
	buf := bytes.NewBuffer([]byte(b))
	for {
		header := next_box_header(buf)
		if header.size == 0 {
			break
		}
		body := next_box_body(buf, header)
		f(body, header)
	}
}

func (this *encoded_box) to_uint32_slice() []uint32 {
	reader := bytes.NewBuffer([]byte(*this))
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	var count uint32
	binary.Read(reader, binary.BigEndian, &count)
	v := make([]uint32, count)
	binary.Read(reader, binary.BigEndian, &v)
	log.Println(count, `to []uint32`)
	return v
}

type full_box_header struct {
	Version byte
	Flags   [3]byte
}

/*
type mp4_box struct {
	ftyp ftyp_box // file type
	moov moov_box // meta container
	mdat mdat_box
}

type moov_box struct {
	mvhd  mvhd_box // movie header
	audio trak_box // audio track/stream
	video trak_box // video track/stream
	//	others []trak_box
}

type trak_box struct {
	tkhd tkhd_box // track header
	//	tref tref_box
	mdia mdia_box // track media container
	// otheres
}
*/

/*
type mdia_box struct {
	mdhd mdhd_box // media header
	hdlr hdlr_box // handler video/audio/hint
	minf minf_box // media information container, not requried
	// video minf box
	// smhd minf box
	// others
}
*/
/*
type minf_box struct {
	vmhd vmhd_box //or smhd gmhd , use this field to determine video/audio track
	hdlr hdlr_box
	dinf dinf_box
	stbl stbl_box
}


type video_minf_box struct {
	vmhd vmhd_box // video media header
	//	smhd smhd_box // sound media heder information
	hdlr hdlr_box
	dinf dinf_box
	stbl stbl_box
}

type sound_minf_box struct {
	smhd smhd_box
	hdlr hdlr_box
	dinf dinf_box
	stbl stbl_box
}
*/

// media info
// contained in mdia
// can contain : gmhd, smhd, vmhd, stbl
/*
type minf_box struct {
	dinf dinf_box //
	stbl stbl_box // sample table
}
type dinf_box struct {
	dref dref_box
}
*/

/*
// data reference
type dref_box struct {
	full_box_header
	count uint32
	data  []dref_box_entry // data references
}

type dref_box_entry struct {
	full_box_header // Flags == 0x0001
	data            []byte
}
*/

/*
// sample table
// contained in minf
// co64, stco, ctts, stsc, stsd, stss, stsz, stts
type stbl_box struct {
	stsd stsd_box // sample table sample descriptions
	stss stss_box // sample table sync sampls/key-frame
	stsz stsz_box // sample table sample size
	stts stts_box // sample table time to sameple map
	stsc stsc_box // sample table sameple to chunk map
	//	stz2 stz2_box // sample size
	stco stco_box // chunk offset
	co64 co64_box // chunk offsets 64bit
	//	ctts ctts_box // 32 bits difference PTS-DS/ composition time offset
}
*/
/*  about ctts
When storing video stream with B-Frames, PTS (Presentation timestamp) may be larger than DTS (Decoder timestamp). It happens because b-frame requires frames following after it do be decoded. Value of this atom is also called Composition Time Offset as, for example, in FLV format.
*/
