package json2xml

// Copyright 2021 Patrick Gundlach

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software
// and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute,
// sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE
// AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT
// OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

// encodeString encodes str as an element named entry with the attribute key="<value of key>" if key is not empty.
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

// ToXML reads the JSON file from jsonInput and writes to xmlOutput.
func ToXML(jsonInput io.Reader, xmlOutput io.Writer) error {
	var err error
	dec := json.NewDecoder(jsonInput)

	enc := xml.NewEncoder(xmlOutput)
	root := xml.StartElement{}
	root.Name = xml.Name{Local: "data"}
	mapElt := xml.StartElement{Name: xml.Name{Local: "map"}}
	aryElt := xml.StartElement{Name: xml.Name{Local: "array"}}
	// inMap is a stack with the last element shows if we are currently
	// in a map. This is to insert key=".." attributes.
	inMap := []bool{false}

	err = enc.EncodeToken(root)
	if err != nil {
		return err
	}

	var key string
encodeLoop:
	for {
		tok, err := dec.Token()
		if err != nil && err != io.EOF {
			return err
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
						return err
					}
					key = ""
				}
				enc.EncodeToken(aryElt.End())
			case '}':
				inMap = inMap[:len(inMap)-1]
				if key != "" {
					err = encodeString(enc, key, "")
					if err != nil {
						return err
					}
					key = ""
				}
				enc.EncodeToken(mapElt.End())
			}
		case string:
			if inMap[len(inMap)-1] {
				if key != "" {
					encodeString(enc, t, key)
					key = ""
				} else {
					key = t
				}

			} else {
				err = encodeString(enc, t, "")
				if err != nil {
					return err
				}
			}
		case float64, bool:
			if key != "" {
				encodeString(enc, fmt.Sprint(t), key)
				key = ""
			} else {
				encodeString(enc, fmt.Sprint(t), "")
			}
		default:
			panic("not implemented")
		}
	}

	enc.EncodeToken(root.End())
	enc.Flush()
	return nil
}
