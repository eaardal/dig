package config

import (
	"time"
)

// Elastic Common Schema (ECS) field names
// https://www.elastic.co/guide/en/ecs/current/ecs-field-reference.html
const (
	ecsMessageField   = "message"
	ecsLevelField     = "log.level"
	ecsTimestampField = "@timestamp"
)

// Logrus field names (copied from logrus package)
const (
	logrusDefaultTimestampFormat = time.RFC3339
	logrusFieldKeyMsg            = "msg"
	logrusFieldKeyLevel          = "level"
	logrusFieldKeyTime           = "time"
	logrusFieldKeyLogrusError    = "logrus_error"
	logrusFieldKeyFunc           = "func"
	logrusFieldKeyFile           = "file"
	logrusErrorKey               = "error"
)

const HomeEnvVar = "DIG_HOME"
const KubernetesNamespaceEnvVar = "DIG_K8S_NAMESPACE"
const KubernetesContextEnvVar = "DIG_K8S_CONTEXT"
const LocalStorageEnvVar = "DIG_LOCAL_STORAGE"

const DigfileFileName = "config.yaml"
const AppName = "dig"

var AppConfig = Config{
	Keywords: KeywordConfig{
		MessageKeywords:   []string{logrusFieldKeyMsg, ecsMessageField},
		LevelKeywords:     []string{logrusFieldKeyLevel, ecsLevelField},
		TimestampKeywords: []string{logrusFieldKeyTime, ecsTimestampField},
		ErrorKeywords:     []string{logrusErrorKey},
		FieldKeywords:     []string{"labels"},
	},
}

type KeywordConfig struct {
	MessageKeywords   []string
	LevelKeywords     []string
	TimestampKeywords []string
	ErrorKeywords     []string
	FieldKeywords     []string
}

type Config struct {
	Keywords KeywordConfig
}
