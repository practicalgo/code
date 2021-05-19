set -e
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt                   -days 365 -subj "/C=AU/ST=NSW/L=Sydney/O=Echorand/OU=Org/CN=practicalgo.echorand.me" -nodes
