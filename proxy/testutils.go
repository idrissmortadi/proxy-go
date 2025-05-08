package proxy

func setupTestConfig() Config {
	return Config{
		Target:     "http://localhost:8080",
		ProxyPort:  8081,
		RateLimit:  1,
		BurstLimit: 1,
		CacheSize:  2,
	}
}
