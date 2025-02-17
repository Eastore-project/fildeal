package buffer

// Config represents the configuration for a buffer
type Config struct {
	Type    string
	ApiKey  string
	BaseURL string
}

// Response represents the response from a buffer operation
type Response struct {
	URL  string
	Hash string
}

// Buffer interface defines the contract for different storage buffers
type Buffer interface {
	Store(filePath string) (*Response, error)
}
