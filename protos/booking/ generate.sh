#!/bin/bash

# Путь к директории с вашими .proto файлами (в вашем случае это внешний репозиторий)
PROTO_DIR=./vendor/github.com/yourusername/booking-proto/proto
# Путь к директории, куда будет генерироваться Go-код
GO_OUT_DIR=./internal/genproto

# Создаем директорию для сгенерированного кода, если ее нет
mkdir -p $GO_OUT_DIR

# Генерируем Go-код из всех .proto файлов в директории booking
protoc -I=$PROTO_DIR \
  --go_out=$GO_OUT_DIR \
  --go_opt=paths=source_relative \
  $PROTO_DIR/booking/*.proto
