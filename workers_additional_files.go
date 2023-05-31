package cloudflare

// Possible script types
const (
	Wasm             = "application/wasm"
	Javascript       = "application/javascript"
	JavascriptModule = "application/javascript+module"
)

type AdditionalFile struct {
	FileName    string
	ScriptType  string
	FileContent string
}
