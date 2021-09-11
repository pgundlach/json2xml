package json2xml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// encodeString encodes an element named string with the attribute key="<value of key>" if key is not empty.
func encodeString(enc *xml.Encoder, str string, key string) error {
	stringElt := xml.StartElement{Name: xml.Name{Local: "entry"}}
	var err error
	if key != "" {
		stringElt.Attr = []xml.Attr{{Name: xml.Name{Local: "key"}, Value: key}}
	}
	if err = enc.EncodeToken(stringElt); err != nil {
		return err
	}
	if err = enc.EncodeToken(xml.CharData([]byte(str))); err != nil {
		return err
	}
	return enc.EncodeToken(stringElt.End())
}

// ToXML returns a string and an error, perhaps nil, from the json encoded in the argument.
func ToXML(r io.Reader) (string, error) {
	var w strings.Builder
	var err error
	dec := json.NewDecoder(r)

	enc := xml.NewEncoder(&w)
	root := xml.StartElement{}
	root.Name = xml.Name{Local: "data"}
	mapElt := xml.StartElement{Name: xml.Name{Local: "map"}}
	aryElt := xml.StartElement{Name: xml.Name{Local: "array"}}
	inMap := []bool{false}

	err = enc.EncodeToken(root)
	if err != nil {
		return "", err
	}

	var key string
encodeLoop:
	for {
		tok, err := dec.Token()
		if err != nil && err != io.EOF {
			return "", err
		}
		if tok == nil {
			break encodeLoop
		}
		switch t := tok.(type) {
		case json.Delim:
			switch t {
			case '{':
				inMap = append(inMap, true)
				if key == "" {
					enc.EncodeToken(mapElt)
				} else {
					attr := xml.Attr{Name: xml.Name{Local: "key"}, Value: key}
					se := mapElt.Copy()
					se.Attr = []xml.Attr{attr}
					err = enc.EncodeToken(se)
				}
				key = ""
			case '[':
				inMap = append(inMap, false)
				if key == "" {
					enc.EncodeToken(aryElt)
				} else {
					attr := xml.Attr{Name: xml.Name{Local: "key"}, Value: key}
					se := aryElt.Copy()
					se.Attr = []xml.Attr{attr}
					err = enc.EncodeToken(se)
				}
				key = ""
			case ']':
				inMap = inMap[:len(inMap)-1]
				if key != "" {
					err = encodeString(enc, key, "")
					if err != nil {
						return "", err
					}
					key = ""
				}
				enc.EncodeToken(aryElt.End())
			case '}':
				inMap = inMap[:len(inMap)-1]
				if key != "" {
					err = encodeString(enc, key, "")
					if err != nil {
						return "", err
					}
					key = ""
				}
				enc.EncodeToken(mapElt.End())
			}
		case string:
			if inMap[len(inMap)-1] {
				if key != "" {
					// string = string
					encodeString(enc, t, key)
					key = ""
				} else {
					key = t
				}

			} else {
				err = encodeString(enc, t, "")
				if err != nil {
					return "", err
				}
			}
		case float64, bool:
			encodeString(enc, fmt.Sprint(t), "")
		default:
			fmt.Println(t)
		}
	}

	enc.EncodeToken(root.End())
	enc.Flush()
	return w.String(), nil
}
