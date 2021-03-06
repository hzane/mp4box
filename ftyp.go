package mp4box

import (
	"bytes"
	"encoding/binary"
)

type ftyp_box struct {
	major_brand       [4]byte
	minor_version     uint32
	compatible_brands [][4]byte
}

func (this ftyp_box) major() string {
	return string(this.major_brand[:])
}
func (this ftyp_box) brands() []string {
	if this.compatible_brands == nil {
		return nil
	}
	brands := make([]string, len(this.compatible_brands))
	for idx, cb := range this.compatible_brands {
		brands[idx] = string(cb[:])
	}
	return brands
}

func (this encoded_box) to_ftyp() (r ftyp_box) {
	var compatible_brands_offset_in_ftyp = 8 // sizeof(major) + sizeof(minor)
	buf := bytes.NewBuffer(this)
	binary.Read(buf, binary.BigEndian, &r.major_brand)
	binary.Read(buf, binary.BigEndian, &r.minor_version)
	cb := (len(this) - compatible_brands_offset_in_ftyp) / 4
	r.compatible_brands = make([][4]byte, cb)
	binary.Read(buf, binary.BigEndian, &r.compatible_brands)
	return
}
