package config

import (
	"os"
	"strconv"
	"time"
)

var config *Config

type Config struct {
	Hetzner HetznerConfig
	Domain string
}


type HetznerConfig struct {
	Token string
	LoadbalanceID int
	LoadBalancerListenerPort int
	LoadBalancerTargetPort int
}




func Get() *Config {
	return config
}

func Init() {
	config = &Config{
		Domain: GetString("DOMAIN", ""),
		Hetzner: HetznerConfig{
			Token: GetString("HETZNER_TOKEN", ""),
			LoadbalanceID: GetInt("HETZNER_LB_ID", 0),
			LoadBalancerListenerPort: GetInt("HETZNER_LB_LISTEN_PORT", 80),
			LoadBalancerTargetPort: GetInt("HETZNER_LB_TARGET_PORT", 80),
		},
	}
}


func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)

	if err != nil {
		return fallback
	}

	return valAsInt
}


func GetDuration(key string, unit time.Duration , fallback int) time.Duration {
	val, ok := os.LookupEnv(key)

	if !ok {
		return unit * time.Duration(fallback)
	}

	valAsInt, err := strconv.Atoi(val)

	if err != nil {
		return unit * time.Duration(fallback)
	}

	return unit * time.Duration(valAsInt)
}


func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	valAsBool, err := strconv.ParseBool(val)

	if err != nil {
		return fallback
	}

	return valAsBool
}