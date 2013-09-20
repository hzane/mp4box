package mp4box

import (
	"bytes"
	"encoding/binary"
)

// sample table sample descriptions
type stsd_box struct {
	//	full_box_header
	count   uint32
	entries []stsd_entry
}

type stsd_entry struct {
	typ  [4]byte
	body encoded_box // headed by sample_entry
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

// generic sample description
type sample_entry struct { // isn't full box
	_ [6]byte // reserved         [6]byte
	_ uint16  // data_refer_index uint16 // dref index
}

type hint_sample_entry struct { // when handler type is hint
	sample_entry
	data []byte
}

type visual_sample_entry struct { // when handler_type is vide
	sample_entry
	_              uint16    // pre_defined     uint16
	_              uint16    // reserved        uint16
	_              [3]uint32 // pre_defined2    [3]uint32
	Width          uint16
	Height         uint16
	HoriResolution uint32 //uint16.uint16
	VertResolution uint32 // uint16.uint16 fixed float
	_              uint32 // reserved2       uint32
	FrameCount     uint16 // 1
	CompressorName [32]byte
	Depth          uint16
	_              int16 // pre_defined3    int16
}

type audio_sample_entry struct { // when handler_type is soun
	sample_entry
	_            [2]uint32 // reserved      [2]uint32
	ChannelCount uint16    // 2
	SampleSize   uint16    // 16
	_            uint16    // pre_defined   uint16
	_            uint16    // reserved2     uint16
	SampleRate   uint32    // time_scale of media << 16
}

/*
type mp4v_box struct {
	visual_sample_entry
	esds esds_box
}

type avc1_box struct {
	visual_sample_entry
	config avcc_box
	// mpeg4bitratebox
	// mpeg4extensiondescriptorbox
}

type mp4a_box struct {
	audio_sample_entry
	esds esds_box
}
*/
/*
type video_sample_description_box struct {
	sample_entry
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
}

// avc decoder configuration
type avcc_box struct { // isn't full box
	AVCDecoderConfigurationRecord []byte
}

type AVCDecoderConfigurationRecord struct {
	configurationVersion       byte // 1
	AVCProfileIndication       byte
	profile_compatibility      byte
	AVCLevelIndication         byte
	lengthSizeMinusOne         byte //111111xx
	numOfSequenceParameterSets byte // 111xxxxx
		[]type struct{  //SPS
			sequenceParameterSetLength uint16
			sequenceParameterSetNALUnit []byte //size is sequenceParameterSetLength
		}  // size is numofsequenceparametersets
		numOfPictureParameterSets byte
		[]type struct {  // PPS
			pictureParameterSetLength uint16
			pictureParameterSetNALUint []byte  // size is pictureParameterSetLength
		}

}
*/
/*
type esds_box struct { // esds This atom contains an MPEG-4 elementary stream descriptor atom, when codec is mp4v
	full_box_header
	ElementaryStreamDescriptor []byte
}


type ElementaryStreamDescriptor struct {
	esid                   uint32
	stream_dependency_flag uint32
	url_flag               byte
	sl_config_descriptor   uint32
	ocr_stream_flag        byte
}*/

/*
type colr_box struct { //colr
	ParameterType         uint32 // `nclc`
	PrimariesIndex        uint16 // 1
	TransferFunctionIndex uint16 // 1
	MatrixIndex           uint16 // 1
}
*/
