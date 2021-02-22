#!/bin/sh
. "${0%/*}/common.sh"

add_mail current/cur/1:2, archive/new/1:2, expected/cur/1:2, <<-EOF
	From: Hans Acker <hans@example.com>
	Subject: Moved from new/ to cur/
	Date: Thu, 23 Mar 2023 15:42:23 +0200
	Message-Id: <EOH1F3NUOY.2KBVMHSBFATNY@example.org>

	This message has been moved from new/ to cur/.
EOF

run_test
