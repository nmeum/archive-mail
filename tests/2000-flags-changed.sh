#!/bin/sh
. "${0%/*}/common.sh"

add_mail current/cur/1:2,S archive/cur/1:2, expected/cur/1:2,S <<-EOF
	From: Hans Acker <hans@example.com>
	Subject: Test
	Date: Thu, 23 Mar 2023 15:42:23 +0200
	Message-Id: <EOH1F3NUOY.2KBVMHSBFATNY@example.org>

	This message has modified flags.
EOF

add_mail current/cur/2:2, archive/cur/2:2, expected/cur/2:2, <<-EOF
	From: Max Mustermann <max@example.com>
	Subject: Unmodified
	Date: Sun, 02 Feb 2020 02:02:02 +0100
	Message-Id: <EOH1F3NUOY.232420SBFATNY@example.org>

	This message is unmodified.
EOF

run_test
