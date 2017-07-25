#!/usr/bin/env bash
# > NB!
# > In Travis this script also tags releases for master branch pushes.
# > This kind of coupling is tolerable, because for us RPM == release.
#
set -e

DESCRIPTION="Mesos 2 IAM"
MAINTAINER=vicent.soria@schibsted.com
URL=http://schibsted.com/
PACKAGE_NAME=mesos2iam-dev

if [[ -n "${TRAVIS_TAG}" ]]; then
    echo "We are building a tag. Skipping..."
    exit 0
fi

# Inside Travis:
# * set dummy version for PR builds
# * tag release for merge pushes
if [[ ${TRAVIS_BRANCH} ]]; then
  VERSION=${TRAVIS_BUILD_NUMBER}
  if [[ "${TRAVIS_BRANCH}" == "master" ]]; then
    git pull --tags
    PREVIOUS_VERSION=$(git tag | grep -i "^v" | sed "s/^v//i" | sort -n | tail -n 1)
    let VERSION=${PREVIOUS_VERSION}+1
    git tag v${VERSION}
    git push --tags
    PACKAGE_NAME=mesos2iam
  fi
fi

# Set dummy version, if not set already (e.g. outside of Travis)
if [[ ! ${VERSION} ]]; then
  VERSION=0.$(date +%Y%m%d%H%M%S)
  echo "WARN: Setting dummy VERSION to: ${VERSION}"
fi

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
PKG_BUILD_DIR="/tmp/rpm.${RANDOM}"; mkdir "${PKG_BUILD_DIR}"

mkdir -p ${PKG_BUILD_DIR}/etc/init ${PKG_BUILD_DIR}/opt/mesos2iam/sbin
cp build/mesos2iam ${PKG_BUILD_DIR}/opt/mesos2iam/sbin/
cp -R ${SCRIPT_DIR}/../deploy/* ${PKG_BUILD_DIR}/
chmod a+x ${PKG_BUILD_DIR}/usr/local/bin/*
pushd build
fpm \
  -s dir \
  -t rpm \
  -n ${PACKAGE_NAME} \
  -v ${VERSION} \
  -m ${MAINTAINER} \
  -d iptables \
  -d sudo \
  --url ${URL} \
  --iteration=$(git rev-parse --short HEAD) \
  --description "${DESCRIPTION}" \
  -C ${PKG_BUILD_DIR}

rm -rf ${PKG_BUILD_DIR}
popd
