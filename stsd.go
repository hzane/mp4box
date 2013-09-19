package mp4box

import (
	"bytes"
	"encoding/binary"
)

// sample table sample descriptions
type stsd_box struct {
	//	full_box_header
	count   int32
	entries []stsd_entry // headed by generic_sample_description
}
type stsd_entry struct {
	typ  [4]byte
	body encoded_box
}

func (this *encoded_box) to_stsd() stsd_box {
	reader := bytes.NewBuffer([]byte(*this))
	//	var h full_box_header
	binary.Read(reader, binary.BigEndian, &full_box_header{})
	var v stsd_box
	binary.Read(reader, binary.BigEndian, &v.count)
	v.entries = make([]stsd_entry, v.count)

	for i := 0; i < int(v.count); i++ {
		h := next_box_header(reader)
		copy(v.entries[i].typ[:], h.typ[:])
		v.entries[i].body = next_box_body(reader, h)
	}
	return v
}

type generic_sample_description struct {
	Reserved       [6]byte
	DataReferIndex uint16 // dref index
}

var (
	sample_description_data_formats = [...]string{`jpeg`, `png `, `mp4v`, `avc1`, `gif `,
		`h263`, `tiff`, `mp4a`, `avcc`}
)

type video_sample_description_box struct {
	generic_sample_description
	Version              uint16   // indicating the version number of the compressed data. This is set to 0
	RevisionLevel        uint16   // 0
	Vendor               uint32   // appl ?
	TemporalQuality      uint32   // [0-1023] the degree of temporal compression.
	SpatialQuality       uint32   // [0-1023] degree of spatial compression.
	Width                uint16   // the width of the source image in pixels
	Height               uint16   // height
	HorizontalResolution uint32   //pixels per inch
	VerticalResolution   uint32   // pixels per inch
	DataSize             uint32   // must be zero
	FrameCount           uint16   // how many frame stored in a sample, usually set to 1
	CompressorName       [32]byte // pascal string
	Depth                uint16   //1,2,4,8,16, 24, 32, 34, 36, 40 color depth. 34,36,40 means 2-,4-,8-bit grayscale
	ColorTableID         int16    // -1 means use default color table. mac color table, 0 means self descripted color table
	//	colr                 colr_box // may be null
	esds esds_box
	avcc avcc_box
}

// avc decoder configuration
type avcc_box struct {
	AVCDecoderConfigurationRecord []byte
}

// only used for 'raw' or 'twos'
type sound_sample_description_v0_box struct {
	generic_sample_description
	version         uint16 // sample description version 0
	revision_level  uint16 // 0
	vendor          uint32 // 0
	number_channels uint16 // 1 mono, 2 stereo, or more
	sample_size     uint16 // 8 = 8 bit , 16 = 16 bit sample
	compression_id  int16  // 0 or -2
	packet_size     uint16 // 0
	sample_rate     uint32 // 16.16 A 32-bit unsigned fixed-point number that indicates the rate at which the sound samples were obtained. This number should match the media’s time scale, that is, the integer portion should match.
	esds            esds_box
}

type sound_sample_description_v1_box struct {
	sound_sample_description_v0_box // version == 1
	samples_per_packet              uint32
	bytes_per_packet                uint32
	bytes_per_frame                 uint32
	bytes_per_sample                uint32
}

// used for mp4a and other
type sound_sample_description_v2_box struct {
	generic_sample_description
	version                            uint16 // ==2
	revision_level                     uint16 // 0
	vendor                             uint32 // 0
	always3                            uint16 // == 3
	always16                           uint16 // 0x0010
	alwaysminus2                       uint16 // -2
	always0                            uint16 // 0
	always65536                        uint32 // 65536
	size_of_struct_only                uint32 // offset to sound sample description
	audio_sample_rate                  float64
	num_audio_channels                 uint32
	always7F000000                     uint32
	const_bits_per_channel             uint32 // only for const or uncompressed audio
	format_specific_flags              uint32 // for lpcm flag
	const_bytes_per_audio_packet       uint32
	const_lpcm_frames_per_audio_packet uint32
}
type esds_box struct { // esds This atom contains an MPEG-4 elementary stream descriptor atom, when codec is mp4v
	full_box_header
	ElementaryStreamDescriptor []byte
}

/*
type colr_box struct { //colr
	ParameterType         uint32 // `nclc`
	PrimariesIndex        uint16 // 1
	TransferFunctionIndex uint16 // 1
	MatrixIndex           uint16 // 1
}
*/
