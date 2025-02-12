# Renew SSL Certs
This script creates a loadbalancer service (Hetzner), starts up a fileserver which meet certbot requirements, checks for health, starts certbot in docker and cleans up after everything is done

## Prerequisite
Following ENV Variables must be set:
```bash
export HETZNER_TOKEN
export HETZNER_LB_ID
export HETZNER_LB_LISTEN_PORT
export HETZNER_LB_TARGET_PORT
export DOMAIN
export EMAIL
```


## Run Script

```bash
go run main.go
```

## New Certs location
The certbot container is mounted to /etc/letsencrypt/ -> Here you will find your new certs