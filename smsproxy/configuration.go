package smsproxy

type smsProxyConfig struct {
	minimumInBatch int
	maxAttempts    int
}

type ConfigOption = func(*smsProxyConfig)

func newConfig() smsProxyConfig {
	return smsProxyConfig{minimumInBatch: 10, maxAttempts: 1}
}

func (config smsProxyConfig) setMinimumInBatch(count int) smsProxyConfig {
	MinimumInBatchOption(count)(&config)
	return config
}

func (config smsProxyConfig) setMaxAttempts(maxAttempts int) smsProxyConfig {
	MaxAttemptsCountOption(maxAttempts)(&config)
	return config
}

func (config smsProxyConfig) disableBatching() smsProxyConfig {
	MinimumInBatchOption(0)(&config)
	return config
}

func MaxAttemptsCountOption(count int) ConfigOption {
	return func(config *smsProxyConfig) {
		config.maxAttempts = count
	}
}

func DisableBatching() ConfigOption {
	return func(config *smsProxyConfig) {
		config.minimumInBatch = 0
	}
}

func MinimumInBatchOption(count int) ConfigOption {
	return func(config *smsProxyConfig) {
		config.minimumInBatch = count
	}
}
