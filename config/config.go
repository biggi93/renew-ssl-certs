package config

import (
	"os"
	"strconv"
	"time"
)

var config *Config

type Config struct {
	Hetzner HetznerConfig
	Domain  string
}

type HetznerConfig struct {
	Token                    string
	LoadbalanceID            int
	FireWallName             string
	LoadBalancerListenerPort int
	LoadBalancerTargetPort   int
	PortSshIn					 string
	PortBWIn						 string
	PortBWOut					 string
	PortCbIn						 string
	PortCbOut					 string
	PortBwMailOut					 string
	Port443					 string
}

func Get() *Config {
	return config
}

func Init() {
	config = &Config{
		Domain: GetString("DOMAIN", ""),
		Hetzner: HetznerConfig{
			Token:                    GetString("HETZNER_TOKEN", ""),
			FireWallName:             GetString("HETZNER_FW_NAME", ""),
			LoadbalanceID:            GetInt("HETZNER_LB_ID", 0),
			LoadBalancerListenerPort: GetInt("HETZNER_LB_LISTEN_PORT", 80),
			LoadBalancerTargetPort:   GetInt("HETZNER_LB_TARGET_PORT", 80),
			PortSshIn: GetString("PORT_SSH_IN", ""),
			PortBWIn: GetString("PORT_BW_IN", ""),
			PortBWOut: GetString("PORT_BW_OUT", ""),
			PortCbIn: GetString("PORT_CERTBOT_IN", ""),
			PortCbOut: GetString("PORT_CERTBOT_OUT", ""),
			PortBwMailOut: GetString("PORT_BW_MAIL_OUT", ""),
			Port443: GetString("PORT_443", ""),
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

func GetDuration(key string, unit time.Duration, fallback int) time.Duration {
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
