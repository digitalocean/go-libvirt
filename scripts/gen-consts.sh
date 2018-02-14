#!/bin/bash
#
# This runs the first code generator used by go-libvirt: c-for-go. This script
# is run from the 'go generate ./...' command, and only needs to be run when
# changing to a different version of libvirt.
if [ -z "${LIBVIRT_SOURCE}" ]; then
    echo "Set LIBVIRT_SOURCE to the root of the libvirt sources you want to use first."
    exit 1
fi

# Make sure c-for-go is installed
if ! which c-for-go > /dev/null; then
    echo "c-for-go not found. Attempting to install it..."
    if ! go get github.com/xlab/c-for-go/...; then
        echo "failed to install c-for-go. Please install it manually from https://github.com/xlab/c-for-go"
        exit 1
    fi
    # HEAD is broken for us because of a failed attempt to redefine the
    # __GNUC_PREREQ macro. Checkout an older version to avoid this.
    pushd $GOPATH/src/github.com/xlab/c-for-go
    git checkout a831fa1954c29f84f7bccab225c3c620d0e7f132
    popd
fi

# Make sure goyacc is installed (needed for the lvgen/ generator)
if ! which goyacc > /dev/null; then
    echo "goyacc not found. Attempting to install it..."
    if ! go get golang.org/x/tools/cmd/goyacc/...; then
        echo "failed to install goyacc. Please install it manually from https://golang.org/x/tools/cmd/goyacc"
        exit 1
    fi
fi

# Temporarily symlink the libvirt sources to a subdirectory because c-for-go
# lacks a mechanism for us to pass it a search path for header files.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LVDIR=lv_source
ln -sF ${LIBVIRT_SOURCE} ${LVDIR}
if ! c-for-go -nostamp -nocgo -ccincl libvirt.yml; then
    echo "c-for-go failed"
    exit 1
fi
mv libvirt/const.go ${DIR}/../const.gen.go
rm ${LVDIR}
rm -rf libvirt/
