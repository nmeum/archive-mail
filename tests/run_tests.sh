#!/bin/sh

cd "${0%/*}"
. common.sh

export ARCHIVE_MAIL="${ARCHIVE_MAIL:-$(pwd)/../archive-mail}"
if [ ! -x "${ARCHIVE_MAIL}" ]; then
	printf "'%s' does not exist\n" "${ARCHIVE_MAIL}" 1>&2
	exit 1
fi

export TESTDIR="${TMPDIR:-/tmp}/mail-archive-tests"
trap "rm -rf '${TESTDIR}' 2>/dev/null" INT EXIT

scriptdir="$(pwd)"
for test in [0-9][0-9][0-9][0-9]:*; do
	mkdir -p "${TESTDIR}"
	cd "${TESTDIR}"
	mmkdir current archive expected

	name="${test##*/}"
	printf "Running test case '%s': " "${name}"

	"${scriptdir}/${test##*/}"
	if [ $? -ne 0 ]; then
		printf "FAIL\n"
		exit 1
	fi

	printf "OK\n"
	cd "${scriptdir}" && rm -rf "${TESTDIR}"
done

exit 0
