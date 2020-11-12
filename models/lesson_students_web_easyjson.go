// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6fabcffeDecodeGithubComTechnoZoomersMasterHubBackendModels(in *jlexer.Lexer, out *LessonStudents) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(LessonStudents, 0, 8)
			} else {
				*out = LessonStudents{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 int64
			v1 = int64(in.Int64())
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6fabcffeEncodeGithubComTechnoZoomersMasterHubBackendModels(out *jwriter.Writer, in LessonStudents) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			out.Int64(int64(v3))
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v LessonStudents) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6fabcffeEncodeGithubComTechnoZoomersMasterHubBackendModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LessonStudents) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6fabcffeEncodeGithubComTechnoZoomersMasterHubBackendModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LessonStudents) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6fabcffeDecodeGithubComTechnoZoomersMasterHubBackendModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LessonStudents) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6fabcffeDecodeGithubComTechnoZoomersMasterHubBackendModels(l, v)
}