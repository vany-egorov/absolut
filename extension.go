package absolut

import "fmt"

type Extension uint8

const (
	UNKNOWN Extension = iota
	JSON
	HTML
	CSV
	TSV
	XML
	YAML
	MSG
	PROTO
)

var extensionText1 = map[Extension]string{
	UNKNOWN: "UNKNOWN",
	JSON:    "JSON",
	HTML:    "HTML",
	CSV:     "CSV",
	TSV:     "TSV",
	XML:     "XML",
	YAML:    "YAML",
	MSG:     "MSG",
	PROTO:   "ProtocolBuffers",
}

var extensionText2 = map[Extension]string{
	UNKNOWN: "unk",
	JSON:    "json",
	HTML:    "html",
	CSV:     "csv",
	TSV:     "tsv",
	XML:     "xml",
	YAML:    "yml",
	MSG:     "msg",
	PROTO:   "proto",
}

var extensionContentType = map[Extension]string{
	UNKNOWN: "application/x-unknown",
	JSON:    "application/json",
	HTML:    "text/html",
	CSV:     "text/csv",
	TSV:     "text/tsv",
	XML:     "text/xml",
	YAML:    "application/x-yaml",
	MSG:     "application/x-msgpack",
	PROTO:   "application/x-protobuf",
}

func ExtensionText(ext Extension) string        { return extensionText1[ext] }
func ExtensionContentType(ext Extension) string { return extensionContentType[ext] }

func (it Extension) Is(that Extension) bool { return it == that }
func (it Extension) IsJSON() bool           { return it.Is(JSON) }
func (it Extension) IsHTML() bool           { return it.Is(HTML) }
func (it Extension) IsCSV() bool            { return it.Is(CSV) }
func (it Extension) IsTSV() bool            { return it.Is(TSV) }
func (it Extension) IsXML() bool            { return it.Is(XML) }
func (it Extension) IsYAML() bool           { return it.Is(YAML) }
func (it Extension) IsMSG() bool            { return it.Is(MSG) }
func (it Extension) IsProto() bool          { return it.Is(PROTO) }
func (it Extension) IsUnknown() bool        { return it.Is(UNKNOWN) }
func (it Extension) GetContentType() string { return ExtensionContentType(it) }
func (it Extension) ContentType() string    { return it.GetContentType() }
func (it Extension) String() string         { return extensionText2[it] }

func (it *Extension) UnmarshalYAML(unmarshal func(interface{}) error) error {
	v := ""
	unmarshal(&v)
	*self = NewEnvironment(v)
	if self.IsUnknown() {
		return fmt.Errorf("got unknown extension '%s'", v)
	}
	return nil
}

func NewExtension(s string) Extension {
	switch s {
	case ".json", ".JSON", "json", "JSON":
		return JSON
	case ".html", ".HTML", "html", "HTML":
		return HTML
	case ".csv", ".CSV", "csv", "CSV":
		return CSV
	case ".tsv", ".TSV", "tsv", "TSV":
		return TSV
	case ".xml", ".XML", "xml", "XML":
		return XML
	case ".yaml", ".yml", "yaml", "yml", "YAML", "YML":
		return YAML
	case ".msg", ".msgpack", "msg", "msgpack", "MSG", "MSGPACK":
		return YAML
	case ".proto", ".PROTO", "proto", "PROTO":
		return PROTO
	}

	return UNKNOWN
}

func GuessExtension(s string) Extension {
	it := NewExtension(s)
	if it.IsUnknown() {
		it = JSON
	}
	return it
}
