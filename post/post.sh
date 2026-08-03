rm -f /opt/bitwarden/bwdata/ssl/pkuepper.de/* &&
cp /opt/renew-ssl-certs/static/live/pkuepper.de/* /opt/bitwarden/bwdata/ssl/pkuepper.de/ &&
chown bitwarden:bitwarden /opt/bitwarden/bwdata/ssl/pkuepper.de/* &&
su bitwarden -c "
cd /opt/bitwarden &&
./bitwarden.sh rebuild &&
./bitwarden.sh restart
"