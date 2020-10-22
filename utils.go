package vlc

import (
	"encoding/xml"
	"fmt"
)

type boolean bool

func (bit *boolean) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "current" {
		*bit = true
	} else if attr.Value == "" {
		*bit = false
	} else {
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", attr.Value)
	}
	return nil
}

func (bit boolean) toBool() bool {
	return bit == true //nolint
}
