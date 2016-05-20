#!/bin/bash

# Verify that all files are correctly golint'd.
EXIT=0
GOLINT=$(golint ./...)

if [[ ! -z $GOLINT ]]; then
	echo $GOLINT
	EXIT=1
fi

exit $EXIT
