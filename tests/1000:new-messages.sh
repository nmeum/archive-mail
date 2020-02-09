#!/bin/sh
. "${0%/*}/common.sh"

add_mail current/cur/1:2, expected/cur/1:2, <<-EOF
	From: Jürgen Jörgsen <juergen@example.com>
	Subject: New Message
	Date: Thu, 23 Mar 2023 15:42:23 +0200
	Message-Id: <EOH1F3NUOY.2KBVMHSBFATNY@example.org>

	This message is new, i.e. not in the archive yet.
EOF

run_test
