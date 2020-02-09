#!/bin/sh
. "${0%/*}/common.sh"

add_mail current/cur/1:2, expected/cur/1:2, expected/cur/1:2, <<-EOF
	From: Some One <example@example.com>
	Subject: Unmodified Message v1
	Date: Mon, 23 Dez 2313 12:23:42 +0200
	Message-Id: <RADNE23UOY.2KBVMHSBFATNY@example.org>

	This message is already in the archive.
EOF

add_mail current/cur/2:2, expected/cur/2:2, expected/cur/2:2, <<-EOF
	From: Some One <example@example.com>
	Subject: Unmodified Message v2
	Date: Mon, 23 Dez 2314 12:23:42 +0200
	Message-Id: <RADNE23UOY.2KBV29410ATNY@example.org>

	This message is already in the archive.
EOF

run_test
