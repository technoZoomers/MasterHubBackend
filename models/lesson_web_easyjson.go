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

func easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels(in *jlexer.Lexer, out *Lessons) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Lessons, 0, 0)
			} else {
				*out = Lessons{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Lesson
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
func easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels(out *jwriter.Writer, in Lessons) {
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
func (v Lessons) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Lessons) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Lessons) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Lessons) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels(l, v)
}
func easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels1(in *jlexer.Lexer, out *Lesson) {
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
		case "id":
			out.Id = int64(in.Int64())
		case "master_id":
			out.MasterId = int64(in.Int64())
		case "time_start":
			out.TimeStart = string(in.String())
		case "time_end":
			out.TimeEnd = string(in.String())
		case "duration":
			out.Duration = string(in.String())
		case "date":
			out.Date = string(in.String())
		case "education_format":
			out.EducationFormat = string(in.String())
		case "price":
			easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels2(in, &out.Price)
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
func easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels1(out *jwriter.Writer, in Lesson) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"master_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.MasterId))
	}
	{
		const prefix string = ",\"time_start\":"
		out.RawString(prefix)
		out.String(string(in.TimeStart))
	}
	{
		const prefix string = ",\"time_end\":"
		out.RawString(prefix)
		out.String(string(in.TimeEnd))
	}
	{
		const prefix string = ",\"duration\":"
		out.RawString(prefix)
		out.String(string(in.Duration))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.String(string(in.Date))
	}
	{
		const prefix string = ",\"education_format\":"
		out.RawString(prefix)
		out.String(string(in.EducationFormat))
	}
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix)
		easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels2(out, in.Price)
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.Int64(int64(in.Status))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Lesson) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Lesson) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Lesson) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Lesson) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels1(l, v)
}
func easyjsonC066239bDecodeGithubComTechnoZoomersMasterHubBackendModels2(in *jlexer.Lexer, out *Price) {
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
		case "value":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Value).UnmarshalJSON(data))
			}
		case "currency":
			out.Currency = string(in.String())
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
func easyjsonC066239bEncodeGithubComTechnoZoomersMasterHubBackendModels2(out *jwriter.Writer, in Price) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix[1:])
		out.Raw((in.Value).MarshalJSON())
	}
	{
		const prefix string = ",\"currency\":"
		out.RawString(prefix)
		out.String(string(in.Currency))
	}
	out.RawByte('}')
}
