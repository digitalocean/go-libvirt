#!/bin/bash
#
# This runs the first code generator used by go-libvirt: c-for-go. This script
# is run from the 'go generate ./...' command, and only needs to be run when
# changing to a different version of libvirt.

# Set TOP to the root of this repo, and SCRIPTS to the absolute path to the 
# scripts/ directory.
TOP="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
SCRIPTS="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ -z "${LIBVIRT_SOURCE}" ]; then
    echo "Set LIBVIRT_SOURCE to the root of the libvirt sources you want to use first."
    exit 1
fi
if [ -z "$GOBIN" ]; then
    export GOBIN="$TOP/bin"
fi

# Ensure this specific tooling comes first on the PATH
export PATH="$GOBIN:$PATH"

# Make sure c-for-go is installed
echo "Attempting to install c-for-go..."
if ! go install github.com/xlab/c-for-go@v1.1.0 ; then
    echo "failed to install c-for-go."
    exit 1
fi

# Make sure goyacc is installed (needed for the lvgen/ generator)
echo "Attempting to install goyacc..."
if ! go install golang.org/x/tools/cmd/goyacc@v0.18.0; then
    echo "failed to install goyacc. Please install it manually from https://golang.org/x/tools/cmd/goyacc"
    exit 1
fi

# Ensure fresh output is in libvirt/
rm -rf libvirt/

# Temporarily symlink the libvirt sources to a subdirectory because c-for-go
# lacks a mechanism for us to pass it a search path for header files.
LVDIR=lv_source
ln -sF "${LIBVIRT_SOURCE}" "${LVDIR}"
if ! c-for-go -nostamp -nocgo -ccincl libvirt.yml; then
    echo "c-for-go failed"
    exit 1
fi

# Use the generated 'const.go'
mv libvirt/const.go "${SCRIPTS}/../const.gen.go"

# Remove unused generated files
rm -f \
    libvirt/cgo_helpers.go \
    libvirt/doc.go \
    libvirt/types.go

rm "${LVDIR}"
