package absolut

type Extension int

const (
	JSON Extension = iota
	HTML
	CSV
	TSV
	XML
)

var extensionText = map[Extension]string{
	JSON: "JSON",
	HTML: "HTML",
	CSV:  "CSV",
	TSV:  "TSV",
	XML:  "XML",
}

func ExtensionText(ext Extension) string {
	return extensionText[ext]
}

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
	}

	return JSON
}
