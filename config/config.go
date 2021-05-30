package config

import (
	"log"
	"os"

	"github.com/mostafasolati/leviathan/contracts"
	"github.com/spf13/viper"

	// enable consul config provider
	_ "github.com/spf13/viper/remote"
)

var envNames = map[string]contracts.Environment{
	"production":  contracts.Production,
	"staging":     contracts.Staging,
	"development": contracts.Development,
}

type configService struct {
	backend *viper.Viper
}

// NewConfigService creates a new IConfigService.
//
// defined by 'CONSUL_HOST' and 'CONSUL_PORT' environment variables, and reads
//
// In each case, it watches for changes and updates configuration parameters
// accordingly.
func NewConfigService(filename string) contracts.IConfigService {
	backend := viper.New()
	addLocalConfigProvider(filename, backend)
	return &configService{backend: backend}
}

func addLocalConfigProvider(filename string, backend *viper.Viper) {
	backend.SetConfigName(filename)
	backend.SetConfigType("yaml")
	backend.AddConfigPath(".")
	if err := backend.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("cannot find config file")
		} else {
			log.Printf("cannot read config file: %v\n", err)
		}
		return
	}

	// Watch local config forever.
	backend.WatchConfig()
}

// IsDebug implements IConfigService.IsDebug
func (s *configService) IsDebug() bool {
	return s.backend.GetBool("debug")
}

// IsProduction implements IConfigService.IsProduction
func (s *configService) IsProduction() bool {
	return s.Environment() == contracts.Production
}

// Environment implements IConfigService.Environment
func (s *configService) Environment() contracts.Environment {
	envName := s.String("environment")
	if env, ok := envNames[envName]; ok {
		return env
	}
	return contracts.Development
}

// BaseURL implements IConfigService.BaseURL
func (s *configService) BaseURL() string {
	return s.String("base-url")
}

// WebURL implements IConfigService.WebURL
func (s *configService) WebURL() string {
	return s.String("web-url")
}

// StorageDir implements IConfigService.StorageDir
func (s *configService) StorageDir() string {
	if dir := s.String("storage-dir"); dir != "" {
		return dir
	}
	if dir, err := os.Getwd(); err != nil {
		return dir
	}
	return "."
}

// StaticDir implements IConfigService.StaticDir
func (s *configService) StaticDir() string {
	if dir := s.String("static-dir"); dir != "" {
		return dir
	}
	if dir, err := os.Getwd(); err != nil {
		return dir
	}
	return "."
}

// Int implements IConfigService.Int
func (s *configService) Int(key string) int {
	return s.backend.GetInt(key)
}

// Bool implements IConfigService.Bool
func (s *configService) Bool(key string) bool {
	return s.backend.GetBool(key)
}

// String implements IConfigService.String
func (s *configService) String(key string) string {
	return s.backend.GetString(key)
}

// SetString implements IConfigService.SetString
func (s *configService) SetString(key, value string) {
	s.backend.Set(key, value)
}

// Dump implements IConfigService.Dump
func (s *configService) Dump() map[string]string {
	dump := make(map[string]string)
	keys := s.backend.AllKeys()
	for _, key := range keys {
		dump[key] = s.String(key)
	}
	return dump
}
