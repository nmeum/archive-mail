add_mail() {
	mail="$(cat)" # absorb here document
	for path in "$@"; do
		echo "${mail}" > "${path}"
	done
}

run_test() {
	# TODO: ensure that current is modified
	"${ARCHIVE_MAIL}" current→archive

	diffout="$(diff -r archive expected)"
	if [ $? -ne 0 ]; then
		printf "FAIL: Output didn't match.\n\n%s\n" "${diffout}"
		exit 1
	fi
}
