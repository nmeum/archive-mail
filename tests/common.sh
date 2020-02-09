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

check_maildir() {
	diffout="$(diff -r "${1}" "${2}")"
	if [ $? -ne 0 ]; then
		printf "FAIL: Directories differ.\n\n%s\n" "${diffout}"
		exit 1
	fi
}

run_test() {
	current="${1:-current}"
	archive="${2:-archive}"
	expected="${3:-expected}"

	cp -r "${current}" "${current}.bkp"
	"${ARCHIVE_MAIL}" "${current}"→"${archive}"

	check_maildir "${archive}" "${expected}"
	if ! diff -r "${current}.bkp" "${current}" >/dev/null; then
		printf "FAIL: current was modified\n"
		exit 1
	fi
}
