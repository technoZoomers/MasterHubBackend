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

func easyjson3915c3ffDecodeGithubComTechnoZoomersMasterHubBackendModels(in *jlexer.Lexer, out *MasterFull) {
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
		case "user_id":
			out.UserId = int64(in.Int64())
		case "email":
			out.Email = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "fullname":
			out.Fullname = string(in.String())
		case "language":
			if in.IsNull() {
				in.Skip()
				out.Languages = nil
			} else {
				in.Delim('[')
				if out.Languages == nil {
					if !in.IsDelim(']') {
						out.Languages = make([]string, 0, 4)
					} else {
						out.Languages = []string{}
					}
				} else {
					out.Languages = (out.Languages)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Languages = append(out.Languages, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "theme":
			(out.Theme).UnmarshalEasyJSON(in)
		case "description":
			out.Description = string(in.String())
		case "qualification":
			out.Qualification = string(in.String())
		case "education_format":
			if in.IsNull() {
				in.Skip()
				out.EducationFormat = nil
			} else {
				in.Delim('[')
				if out.EducationFormat == nil {
					if !in.IsDelim(']') {
						out.EducationFormat = make([]string, 0, 4)
					} else {
						out.EducationFormat = []string{}
					}
				} else {
					out.EducationFormat = (out.EducationFormat)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.EducationFormat = append(out.EducationFormat, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "avg_price":
			easyjson3915c3ffDecodeGithubComTechnoZoomersMasterHubBackendModels1(in, &out.AveragePrice)
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
func easyjson3915c3ffEncodeGithubComTechnoZoomersMasterHubBackendModels(out *jwriter.Writer, in MasterFull) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.UserId))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"fullname\":"
		out.RawString(prefix)
		out.String(string(in.Fullname))
	}
	{
		const prefix string = ",\"language\":"
		out.RawString(prefix)
		if in.Languages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Languages {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"theme\":"
		out.RawString(prefix)
		(in.Theme).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"qualification\":"
		out.RawString(prefix)
		out.String(string(in.Qualification))
	}
	{
		const prefix string = ",\"education_format\":"
		out.RawString(prefix)
		if in.EducationFormat == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.EducationFormat {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"avg_price\":"
		out.RawString(prefix)
		easyjson3915c3ffEncodeGithubComTechnoZoomersMasterHubBackendModels1(out, in.AveragePrice)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MasterFull) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3915c3ffEncodeGithubComTechnoZoomersMasterHubBackendModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MasterFull) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3915c3ffEncodeGithubComTechnoZoomersMasterHubBackendModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MasterFull) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3915c3ffDecodeGithubComTechnoZoomersMasterHubBackendModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MasterFull) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3915c3ffDecodeGithubComTechnoZoomersMasterHubBackendModels(l, v)
}
func easyjson3915c3ffDecodeGithubComTechnoZoomersMasterHubBackendModels1(in *jlexer.Lexer, out *Price) {
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
func easyjson3915c3ffEncodeGithubComTechnoZoomersMasterHubBackendModels1(out *jwriter.Writer, in Price) {
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