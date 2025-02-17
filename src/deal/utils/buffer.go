package dealutils

// BufferConfig represents the configuration for a buffer
type BufferConfig struct {
	Type    string
	ApiKey  string
	BaseURL string
}

// BufferResponse represents the response from a buffer operation
type BufferResponse struct {
	URL      string
	Hash     string
	Location string
}

// Buffer interface defines the contract for different storage buffers
type Buffer interface {
	Store(filePath string) (*BufferResponse, error)
}

// LocalBuffer implements Buffer for local storage
type LocalBuffer struct{}

func NewLocalBuffer() Buffer {
	return &LocalBuffer{}
}

func (b *LocalBuffer) Store(filePath string) (*BufferResponse, error) {
	return &BufferResponse{
		Location: filePath,
	}, nil
}

// LighthouseBuffer implements Buffer for Lighthouse storage
type LighthouseBuffer struct {
	apiKey  string
	baseURL string
}

func NewLighthouseBuffer(apiKey, baseURL string) Buffer {
	return &LighthouseBuffer{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (b *LighthouseBuffer) Store(filePath string) (*BufferResponse, error) {
	resp, err := UploadToLighthouse(filePath, b.apiKey)
	if err != nil {
		return nil, err
	}
	return &BufferResponse{
		URL:  b.baseURL + resp.Hash,
		Hash: resp.Hash,
	}, nil
}

// NewBuffer creates a new buffer based on the config
func NewBuffer(config *BufferConfig) Buffer {
	switch config.Type {
	case "lighthouse":
		return NewLighthouseBuffer(config.ApiKey, config.BaseURL)
	default:
		return NewLocalBuffer()
	}
}
