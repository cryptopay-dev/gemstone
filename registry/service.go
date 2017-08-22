package registry

type Service struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    int    `json:"port"`
	Addr    string `json:"addr"`
}
