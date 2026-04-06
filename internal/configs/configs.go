package configs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Log struct {
	Level  string `mapstructure:"level" json:"level" validate:"oneof=debug info warn error"`
	Writer string `mapstructure:"writer" json:"writer" validate:"oneof=stdout file"`
	Path   string `mapstructure:"path" json:"path" validate:"required_if=Writer file"`
}

type Config struct {
	filepath            string
	ServiceName         string `mapstructure:"service_name" json:"service_name"`
	Version             string `mapstructure:"version" json:"version"`
	SearchFolder        string `mapstructure:"search_folder" json:"search_folder"`
	DegreeOfParallelism int    `mapstructure:"degree_of_parallelism" json:"degree_of_parallelism" validate:"min=1"`

	Log Log `mapstructure:"log" json:"log"`
}

func (c *Config) SetVersion(version string) {
	c.Version = version
}

func (c *Config) AttachData() error {
	viper.SetDefault("service_name", "mm")
	viper.SetDefault("version", "v0.0.0")
	viper.SetDefault("search_folder", ".")
	viper.SetDefault("degree_of_parallelism", 10)

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.writer", "stdout")

	if c.filepath != "" {
		filepathClean := filepath.Clean(c.filepath)
		filename := filepath.Base(filepathClean)
		extFull := filepath.Ext(filepathClean)
		dirPath := filepath.Dir(filepathClean)
		basename := strings.TrimSuffix(filename, extFull)
		extName := strings.TrimPrefix(extFull, ".")

		viper.SetConfigName(basename)
		viper.SetConfigType(extName)
		viper.AddConfigPath(dirPath)

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(c); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func NewConfig(filepath string) Config {
	return Config{
		filepath: filepath,
	}
}
