#!/bin/bash

go_build() {
  go get ./...
  CGO_ENABLED=0 go build ./...
}

go_build