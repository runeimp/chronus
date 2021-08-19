PROJECT_NAME := "Chronus"
PROJECT_CLI := "chronus"

alias arc := archive

@_default:
	just _term-wipe
	just --list


# Archive GoReleaser dist
archive:
	#!/bin/sh
	just _term-wipe
	tag="$(git tag --points-at main)"
	app="{{PROJECT_NAME}}"
	arc="${app}_${tag}"

	# echo "app = '${app}'"
	# echo "tag = '${tag}'"
	# echo "arc = '${arc}'"
	if [ ! -e distro ]; then
		mkdir distro
	fi
	if [ -e dist ]; then
		echo "Move dist -> distro/${arc}"
		mv dist "distro/${arc}"

		# echo "cd distro"
		cd distro
		
		printf "pwd = "
		pwd
		
		ls -Alh
	else
		echo "dist directory not found for archiving"
	fi


# Build and install app
build:
	@just _term-wipe
	go build -o {{PROJECT_CLI}} ./cmd/{{PROJECT_CLI}}/main.go
	mv {{PROJECT_CLI}} "${GOBIN}/"


# Build distro
distro:
	#!/bin/sh
	goreleaser
	just archive


install:
	#!/bin/sh
	# go install ./cmd/{{PROJECT_CLI}}/main.go
	cd cmd/{{PROJECT_CLI}}
	go install


# Run code
run +args='"2021-03-08 16:06:34 MST"':
	@just _term-wipe
	CHRONUS_COUNTRY_CODE="US" go run ./cmd/{{PROJECT_CLI}}/main.go {{args}}
	@#hr; echo
	@#go run ./cmd/{{PROJECT_CLI}}/main.go -country-code USA {{args}}
	@#hr; echo
	@#go run ./cmd/{{PROJECT_CLI}}/main.go {{args}}
	@#hr; echo
	@#go run ./cmd/{{PROJECT_CLI}}/main.go -h


# Run a test
@test cmd="coverage":
	just _term-wipe
	just test-{{cmd}}

# Run Go Unit Tests
@test-coverage:
	just _term-wipe
	echo "You need to run:"
	echo "go test -coverprofile=c.out"
	echo "go tool cover -func=c.out"


_term-wipe:
	#!/usr/bin/env bash
	set -exo pipefail
	if [[ ${#VISUAL_STUDIO_CODE} -gt 0 ]]; then
		clear
	elif [[ ${KITTY_WINDOW_ID} -gt 0 ]] || [[ ${#TMUX} -gt 0 ]] || [[ "${TERM_PROGRAM}" = 'vscode' ]]; then
		printf '\033c'
	elif [[ "${TERM_PROGRAM}" = 'Apple_Terminal' ]] || [[ "${TERM_PROGRAM}" = 'iTerm.app' ]]; then
		osascript -e 'tell application "System Events" to keystroke "k" using command down'
	elif [[ -x "$(which tput)" ]]; then
		tput reset
	elif [[ -x "$(which tcap)" ]]; then
		tcap rs
	elif [[ -x "$(which reset)" ]]; then
		reset
	else
		clear
	fi

