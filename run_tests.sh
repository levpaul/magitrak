#!/bin/bash
find ./tests/ -name '*_test.go' -exec echo "{}" \; -exec go test "{}" \; ;
