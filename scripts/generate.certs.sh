#!/bin/sh
openssl req -new -subj "/CN=localhost" -newkey rsa:2048 -nodes -keyout ./server.key -out ./server.csr
openssl x509 -req -days 3650 -in ./server.csr -signkey ./server.key -out ./server.crt