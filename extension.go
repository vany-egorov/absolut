package absolut

type Extension uint8

const (
	JSON Extension = iota
	HTML
	CSV
	TSV
	XML
	PROTO
)

var extensionText = map[Extension]string{
	JSON:  "JSON",
	HTML:  "HTML",
	CSV:   "CSV",
	TSV:   "TSV",
	XML:   "XML",
	PROTO: "ProtocolBuffers",
}

var extensionContentType = map[Extension]string{
	JSON:  "application/json",
	HTML:  "text/html",
	CSV:   "text/csv",
	TSV:   "text/tsv",
	XML:   "text/xml",
	PROTO: "application/x-protobuf",
}

func ExtensionText(ext Extension) string        { return extensionText[ext] }
func ExtensionContentType(ext Extension) string { return extensionContentType[ext] }

func (it Extension) Is(that Extension) bool { return it == that }
func (it Extension) IsJSON() bool           { return it.Is(JSON) }
func (it Extension) IsHTML() bool           { return it.Is(HTML) }
func (it Extension) IsCSV() bool            { return it.Is(CSV) }
func (it Extension) IsTSV() bool            { return it.Is(TSV) }
func (it Extension) IsXML() bool            { return it.Is(XML) }
func (it Extension) IsProto() bool          { return it.Is(PROTO) }
func (it Extension) GetContentType() string { return ExtensionContentType(it) }

func GuessExtension(s string) Extension {
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
	case ".proto", ".PROTO", "proto", "PROTO":
		return PROTO
	}

	return JSON
}
