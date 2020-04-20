package servicers

// AppConfig provides the config specific to a download application
type AppConfig struct {
	s3Region    string
	s3Bucket    string
	s3SubFolder string
}

// DownloadServiceConfig provides the set of configs for the download service
type DownloadServiceConfig struct {
	apps map[string]AppConfig
}

func InitServiceConfig() DownloadServiceConfig {
	config := DownloadServiceConfig{}
	config.apps = make(map[string]AppConfig)
	// TODO: Load this from a yaml file
	config.apps["feg"] = AppConfig{"us-west-2", "magma-images", "feg/"}
	config.apps["soma"] = AppConfig{"us-east-1", "soma.images", ""}
	// Add a default service for backward compatibility
	config.apps["default"] = AppConfig{"us-east-1", "soma.images", ""}
	return config
}
