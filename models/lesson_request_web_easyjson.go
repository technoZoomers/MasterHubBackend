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

func easyjson3fe14cc9DecodeGithubComTechnoZoomersMasterHubBackendModels(in *jlexer.Lexer, out *LessonRequests) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(LessonRequests, 0, 2)
			} else {
				*out = LessonRequests{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 LessonRequest
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3fe14cc9EncodeGithubComTechnoZoomersMasterHubBackendModels(out *jwriter.Writer, in LessonRequests) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v LessonRequests) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3fe14cc9EncodeGithubComTechnoZoomersMasterHubBackendModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LessonRequests) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3fe14cc9EncodeGithubComTechnoZoomersMasterHubBackendModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LessonRequests) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3fe14cc9DecodeGithubComTechnoZoomersMasterHubBackendModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LessonRequests) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3fe14cc9DecodeGithubComTechnoZoomersMasterHubBackendModels(l, v)
}
func easyjson3fe14cc9DecodeGithubComTechnoZoomersMasterHubBackendModels1(in *jlexer.Lexer, out *LessonRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "lesson_id":
			out.LessonId = int64(in.Int64())
		case "student_id":
			out.StudentId = int64(in.Int64())
		case "status":
			out.Status = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3fe14cc9EncodeGithubComTechnoZoomersMasterHubBackendModels1(out *jwriter.Writer, in LessonRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"lesson_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.LessonId))
	}
	{
		const prefix string = ",\"student_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.StudentId))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.Int64(int64(in.Status))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LessonRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3fe14cc9EncodeGithubComTechnoZoomersMasterHubBackendModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LessonRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3fe14cc9EncodeGithubComTechnoZoomersMasterHubBackendModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LessonRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3fe14cc9DecodeGithubComTechnoZoomersMasterHubBackendModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LessonRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3fe14cc9DecodeGithubComTechnoZoomersMasterHubBackendModels1(l, v)
}
