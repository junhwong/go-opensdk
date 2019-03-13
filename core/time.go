package core

import (
	"encoding/xml"
	"time"
)

type Time struct {
	time.Time
	Format string
}

//MarshalJSON 实现它的json序列化方法
func (t Time) MarshalJSON() ([]byte, error) {
	f := t.Format
	if f == "" {
		f = "2006-01-02 15:04:05"
	}
	s := t.Time.Format(f)
	return []byte("\"" + s + "\""), nil
}

//MarshalXML 实现它的json序列化方法
func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	f := t.Format
	if f == "" {
		f = "2006-01-02 15:04:05"
	}
	e.EncodeElement(t.Time.Format(f), start)
	// s := t.Time.Format(f)
	return nil
}
