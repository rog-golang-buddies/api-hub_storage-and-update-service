package config

// Web is a web-related properties configuration
type Web struct {
	//RespLimBytes represents the maximum file size (in bytes) to download.
	RespLimBytes int64 `default:"5242880"`
}
