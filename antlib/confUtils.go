package antlib

import (
	"antgo/antlib/config"
	"path/filepath"
	"os"
)

type antConfig struct {
	innerConfig config.Configer
}

var (
	AppConfig         *antConfig
	AppAntPath        string
	appConfigPath     string
	appConfigProvider = "ini"
)

func AntInit(path string) {
	var err error
	if AppAntPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath = filepath.Join(workPath, path)
	if !FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppAntPath, path)
		if !FileExists(appConfigPath) {
			AppConfig = &antConfig{innerConfig: config.NewFakeConfig()}
			return
		}
	}

	if err = parseAntConfig(appConfigPath); err != nil {
		panic(err)
	}

}

// now only support ini, next will support json.
func parseAntConfig(appConfigPath string) (err error) {
	AppConfig, err = newAntConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	return
}
func newAntConfig(appConfigProvider, appConfigPath string) (*antConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &antConfig{ac}, nil
}
func (this *antConfig) Set(key, val string) error {
	return this.innerConfig.Set(key, val)
}
func (this *antConfig) String(key string) string {
	return this.innerConfig.String(key)
}
func (b *antConfig) Strings(key string) []string {
	return b.innerConfig.Strings(key)
}

func (b *antConfig) Int(key string) (int, error) {

	return b.innerConfig.Int(key)
}

func (b *antConfig) Int64(key string) (int64, error) {

	return b.innerConfig.Int64(key)
}

func (b *antConfig) Bool(key string) (bool, error) {

	return b.innerConfig.Bool(key)
}

func (b *antConfig) Float(key string) (float64, error) {

	return b.innerConfig.Float(key)
}

func (b *antConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *antConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *antConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *antConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *antConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *antConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *antConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *antConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *antConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
