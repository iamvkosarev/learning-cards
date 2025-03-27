ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)

-include .env
export
