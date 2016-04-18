#!/bin/bash

SIZE=${1:-512}
openssl genrsa -out app.rsa ${SIZE}
openssl rsa -in app.rsa -pubout > app.rsa.pub
