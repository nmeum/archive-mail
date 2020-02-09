mmkdir() {
	for maildir in "$@"; do
		mkdir -p "${maildir}/new" "${maildir}/cur" "${maildir}/tmp"
	done
}

add_mail() {
	mail="$(cat)" # absorb here document
	for path in "$@"; do
		echo "${mail}" > "${path}"
	done
}

run_test() {
	current="${1:-current}"
	archive="${2:-archive}"
	expected="${3:-expected}"

	cp -r "${current}" "${current}.bkp"
	"${ARCHIVE_MAIL}" "${current}"→"${archive}"

	diffout="$(diff -r "${archive}" "${expected}")"
	if [ $? -ne 0 ]; then
		printf "FAIL: Output didn't match.\n\n%s\n" "${diffout}"
		exit 1
	fi

	if ! diff -r "${current}.bkp" "${current}" >/dev/null; then
		printf "FAIL: current was modified\n"
		exit 1
	fi
}
