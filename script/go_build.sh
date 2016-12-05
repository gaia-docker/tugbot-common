#!/bin/bash

go_build() {
  glide install
  CGO_ENABLED=0 go build ./...
}

go_build