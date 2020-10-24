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

func easyjson83e04af4DecodeGithubComTechnoZoomersMasterHubBackendModels(in *jlexer.Lexer, out *Messages) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Messages, 0, 1)
			} else {
				*out = Messages{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Message
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
func easyjson83e04af4EncodeGithubComTechnoZoomersMasterHubBackendModels(out *jwriter.Writer, in Messages) {
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
func (v Messages) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson83e04af4EncodeGithubComTechnoZoomersMasterHubBackendModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Messages) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson83e04af4EncodeGithubComTechnoZoomersMasterHubBackendModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Messages) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson83e04af4DecodeGithubComTechnoZoomersMasterHubBackendModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Messages) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson83e04af4DecodeGithubComTechnoZoomersMasterHubBackendModels(l, v)
}
func easyjson83e04af4DecodeGithubComTechnoZoomersMasterHubBackendModels1(in *jlexer.Lexer, out *Message) {
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
		case "type":
			out.Type = int64(in.Int64())
		case "author_id":
			out.AuthorId = int64(in.Int64())
		case "text":
			out.Text = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
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
func easyjson83e04af4EncodeGithubComTechnoZoomersMasterHubBackendModels1(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.Int64(int64(in.Type))
	}
	if in.AuthorId != 0 {
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.AuthorId))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson83e04af4EncodeGithubComTechnoZoomersMasterHubBackendModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson83e04af4EncodeGithubComTechnoZoomersMasterHubBackendModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson83e04af4DecodeGithubComTechnoZoomersMasterHubBackendModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson83e04af4DecodeGithubComTechnoZoomersMasterHubBackendModels1(l, v)
}
