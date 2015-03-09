#!/bin/bash

args=${@:-"--short"}

go test ./... $args
