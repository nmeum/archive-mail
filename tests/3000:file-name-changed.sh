#!/bin/sh
. "${0%/*}/common.sh"

add_mail current/cur/new-name:2, archive/cur/old-name:2, expected/cur/new-name:2, <<-EOF
	From: Hans Acker <hans@example.com>
	Subject: Different file name
	Date: Thu, 23 Mar 2023 15:42:23 +0200
	Message-Id: <EOH1F3NUOY.2KBVMHSBFATNY@example.org>

	This message has a new file name.
EOF

run_test
