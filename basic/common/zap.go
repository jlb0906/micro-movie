package common

// Zap 配置
type Zap struct {
	Level       string   `json:"level"`
	Development bool     `json:"development"`
	LogFileDir  string   `json:"logFileDir"`
	OutputPaths []string `json:"outputPaths"`
	MaxSize     int      `json:"maxSize"`
	MaxBackups  int      `json:"maxBackups"`
	MaxAge      int      `json:"maxAge"`
}
