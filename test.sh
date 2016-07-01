#!/bin/bash

export GIN_MODE=release

go test "$@" 
