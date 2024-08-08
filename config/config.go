package config

import (
	"github.com/mateothegreat/go-config/config"
	"github.com/mateothegreat/tailer/config"
)

// Config is the configuration for the tailer.
type Config struct {
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Name    string `yaml:"name" env:"NAME"`
	Address string `yaml:"address" env:"ADDRESS"`
	Port    int    `yaml:"port" env:"PORT"`
}

// GetConfig returns a config of type T.
// It will merge the base config with the environment config.
// If the environment config does not exist, it will use the base config.
//
// Arguments:
//   - env: The environment to use.
//
// Returns:
//   - A pointer to the config of type T.
//   - An error if the config could not be found.
func GetConfig(profiles []string) ([]Profile, error) {
	config, err := config.GetConfig[Config](config.GetConfigArgs{
		Paths: []string{
			".tailor.yaml",
			".tailor.json",
		},
		WalkDepth: 6,
	})
	if err != nil {
		return nil, err
	}

	configs := make([]Profile, len(profiles))
	for i, p := range profiles {
		for _, profile := range config.Profiles {
			if profile.Name == p {
				configs[i] = profile
			}
		}
	}

	return configs, nil
}
