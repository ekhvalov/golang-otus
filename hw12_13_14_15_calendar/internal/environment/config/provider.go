package config

type Provider interface {
	UnmarshalKey(key string, rawVal interface{}) error
	GetString(key string) string
}
