package route

type Option struct {
	Type   uint8       `json:"type"`
	Host   string      `json:"host,omitempty"`
	Listen string      `json:"listen,omitempty"`
	Routes interface{} `json:"routes,omitempty"`
}
