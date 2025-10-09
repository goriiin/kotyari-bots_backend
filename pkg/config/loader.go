package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	env             = ".env"
	appEnv          = "ENV"
	envFile         = ".env"
	envLocal        = "local"
	envProduction   = "prod"
	prodPrefix      = "PROD_"
	localPrefix     = "LOCAL_"
	localPrefixYaml = "local_"
)

var (
	configName  string
	once        sync.Once
	loader      *ConfigLoader
	searchPaths = []string{
		".",
		"./configs",
		"../configs",
	}
)

func getLoader(confName string) *ConfigLoader {
	once.Do(func() {
		env := flag.String("env", "", "Environment: prod, local")
		configPath := flag.String("config", "", "Path to config file")
		flag.Parse()
		configName = confName

		loader = newConfigLoader(*env, *configPath)
	})

	return loader
}

func Load[T any](config *T) error {
	return getLoader(configName).Load(config)
}

func New[T any]() (*T, error) {
	cfg := new(T)
	err := Load(cfg)
	return cfg, err
}

func NewWithConfig[T any](confName string) (*T, error) {
	cfg := new(T)
	configName = confName
	err := Load(cfg)
	return cfg, err
}

func GetViper() *viper.Viper {
	return getLoader(configName).viper
}

type ConfigLoader struct {
	environment string
	viper       *viper.Viper
	mu          sync.RWMutex
	defaults    map[string]interface{}

	envFlag    string
	configPath string
}

func newConfigLoader(envFlag, configPath string) *ConfigLoader {
	cl := &ConfigLoader{
		viper:      viper.New(),
		defaults:   make(map[string]interface{}),
		envFlag:    envFlag,
		configPath: configPath,
	}
	cl.initialize()
	return cl
}

func (cl *ConfigLoader) loadEnvFiles() {
	for _, path := range searchPaths {
		fullPath := filepath.Join(path, envFile)
		if _, err := os.Stat(fullPath); err == nil {
			if err = godotenv.Load(fullPath); err == nil {
				log.Printf("Loaded env file: %s\n", fullPath)
			}
		}
	}

	if cl.configPath != "" {
		dir := filepath.Dir(cl.configPath)
		fullPath := filepath.Join(dir, envFile)
		if _, err := os.Stat(fullPath); err == nil {
			err = godotenv.Load(fullPath)
			if err != nil {
				log.Printf("error when godotenv.Load: %v", err)
				return
			}
		}
	}
}

func (cl *ConfigLoader) initialize() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.loadEnvFiles()
	cl.detectEnvironment(cl.envFlag)

	cl.viper.SetConfigName(configName)
	cl.viper.SetConfigType("yaml")

	if cl.configPath != "" {
		cl.viper.SetConfigFile(cl.configPath)
	} else {
		cl.viper.AddConfigPath(".")
		cl.viper.AddConfigPath("./configs")
	}

	cl.viper.AutomaticEnv()
	cl.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for key, value := range cl.defaults {
		cl.viper.SetDefault(key, value)
	}

	err := cl.viper.ReadInConfig()
	if err != nil {
		log.Printf("error when cl.viper.ReadInConfig(): %v", err)
	}
	cl.loadEnvironmentVariables()
	if cl.environment == envLocal {
		cl.applyLocalOverrides()
	}
}

func (cl *ConfigLoader) loadEnvironmentVariables() {
	var envPrefix string
	switch cl.environment {
	case envLocal:
		envPrefix = localPrefix
	case envProduction:
		envPrefix = prodPrefix
	default:
		return
	}

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]
		if !strings.HasPrefix(key, envPrefix) {
			continue
		}
		viperKey := strings.TrimPrefix(key, envPrefix)

		viperKey = strings.ToLower(viperKey)
		viperKey = strings.ReplaceAll(viperKey, "_", ".")

		cl.viper.Set(viperKey, value)
		if cl.environment == envLocal {
			log.Printf("Set %s = %s", viperKey, value)
		}
	}
}

func (cl *ConfigLoader) detectEnvironment(envFlag string) {
	// Приоритет определения окружения:
	// 1. Флаг командной строки --env
	// 2. Переменная окружения ENV
	// 3. По умолчанию - local

	if envFlag != "" {
		switch envFlag {
		case envProduction:
			cl.environment = envProduction
		case envLocal:
			cl.environment = envLocal
		default:
			cl.environment = envLocal
		}
	} else if env := os.Getenv(appEnv); env != "" {
		switch env {
		case "prod":
			cl.environment = envProduction
		case "local":
			cl.environment = envLocal
		default:
			cl.environment = envLocal
		}
	} else {
		cl.environment = envLocal
	}

	cl.viper.Set("environment", cl.environment)
}

func (cl *ConfigLoader) applyLocalOverrides() {
	settings := cl.viper.AllSettings()

	cl.applyLocalOverridesRecursive(settings, "")

	for key, value := range settings {
		cl.viper.Set(key, value)
	}
}

func (cl *ConfigLoader) applyLocalOverridesRecursive(settings map[string]interface{}, prefix string) {
	keysToDelete := []string{}

	for key, value := range settings {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if strings.HasPrefix(key, localPrefixYaml) {
			originalKey := strings.TrimPrefix(key, localPrefixYaml)
			settings[originalKey] = value
			keysToDelete = append(keysToDelete, key)
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			cl.applyLocalOverridesRecursive(nestedMap, fullKey)
		}
	}

	for _, key := range keysToDelete {
		delete(settings, key)
	}
}

func (cl *ConfigLoader) Load(config interface{}) error {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct")
	}

	if !cl.hasConfigBase(rv.Elem().Type()) {
		return fmt.Errorf("config struct must embed ConfigBase")
	}

	if err := cl.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cl.applyEnvOverrides(config)

	cl.setBaseFields(config)

	if validator, ok := config.(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

func (cl *ConfigLoader) applyEnvOverrides(config interface{}) {
	var envPrefix string
	switch cl.environment {
	case envLocal:
		envPrefix = localPrefix
	case envProduction:
		envPrefix = prodPrefix
	default:
		return
	}

	rv := reflect.ValueOf(config).Elem()
	cl.applyEnvOverridesRecursive(rv, "", envPrefix)
}

func (cl *ConfigLoader) applyEnvOverridesRecursive(v reflect.Value, prefix, envPrefix string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		fieldName := fieldType.Tag.Get("mapstructure")
		if fieldName == "" || fieldName == "-" {
			fieldName = strings.ToLower(fieldType.Name)
		}

		envPath := fieldName
		if prefix != "" {
			envPath = prefix + "_" + fieldName
		}

		if field.Kind() == reflect.Struct && fieldType.Type.Name() != "ConfigBase" {
			cl.applyEnvOverridesRecursive(field, envPath, envPrefix)
		} else {
			envKey := envPrefix + strings.ToUpper(envPath)
			if envValue, exists := os.LookupEnv(envKey); exists && envValue != "" {
				cl.setFieldValue(field, envValue)
			}
		}
	}
}

func (cl *ConfigLoader) setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(uintVal)
		}
	case reflect.Float32, reflect.Float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatVal)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolVal)
		}
	}
}

func (cl *ConfigLoader) hasConfigBase(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Name() == "ConfigBase" && field.Anonymous {
			return true
		}
	}
	return false
}

func (cl *ConfigLoader) setBaseFields(config interface{}) {
	v := reflect.ValueOf(config).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Type().Name() == "ConfigBase" && field.CanSet() {
			base := ConfigBase{
				Environment: cl.environment,
			}
			field.Set(reflect.ValueOf(base))
			break
		}
	}
}

func SetDefault(key string, value interface{}) {
	if loader == nil {
		if defaults == nil {
			defaults = make(map[string]interface{})
		}
		defaults[key] = value

		return
	}

	getLoader(configName).viper.SetDefault(key, value)
}

var defaults map[string]interface{}

type ChangeHandler func()

func WatchConfig(onChange ChangeHandler) {
	l := getLoader(configName)
	l.viper.WatchConfig()
	l.viper.OnConfigChange(func(e fsnotify.Event) {
		l.mu.Lock()
		l.detectEnvironment("")
		l.loadEnvironmentVariables()
		if l.environment == envLocal {
			l.applyLocalOverrides()
		}
		l.mu.Unlock()

		if onChange != nil {
			onChange()
		}
	})
}

func (cl *ConfigLoader) DebugConfig() {
	fmt.Println("=== Current Configuration ===")
	fmt.Printf("Environment: %s\n", cl.environment)
	fmt.Println("\nAll settings:")
	for k, v := range cl.viper.AllSettings() {
		fmt.Printf("%s: %v\n", k, v)
	}
}
