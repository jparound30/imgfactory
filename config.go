package main

type Config struct {
	Port         int
	ImageOptions ImageOptions
	FontOptions  FontOptions
}

type ImageOptions struct {
	Width  int
	Height int
}
type FontOptions struct {
	SizeInPoint float64
}

var config *Config

func LoadConfig() *Config {
	if config != nil {
		return config
	}
	config = &Config{
		Port: 8080,
		ImageOptions: ImageOptions{
			Width:  640,
			Height: 480,
		},
		FontOptions: FontOptions{
			SizeInPoint: 36.0,
		},
	}
	return config
}
