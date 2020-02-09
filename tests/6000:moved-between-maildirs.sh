#!/bin/sh
. "${0%/*}/common.sh"

mmkdir current1 archive1 expected1 \
	current2 archive2 expected2

add_mail current2/cur/1:2, archive1/cur/1:2, expected2/cur/1:2, <<-EOF
	From: Hans Acker <hans@example.com>
	Subject: Moved between maildirs
	Date: Thu, 23 Mar 2023 15:42:23 +0200
	Message-Id: <EOH1F3NUOY.2KBVMHSBFATNY@example.org>

	This message was in maildir1 and moved to maildir2.
EOF

add_mail current1/cur/2:2, archive1/cur/2:2, expected1/cur/2:2, <<-EOF
	From: Some One <example@example.com>
	Subject: Unmodified Message v1
	Date: Mon, 23 Dez 2313 12:23:42 +0200
	Message-Id: <RADNE23UOY.2KBVMHSBFATNY@example.org>

	This message is already in the archive.
EOF

add_mail current2/cur/3:2, expected2/cur/3:2, expected2/cur/3:2, <<-EOF
	From: Some One <example@example.com>
	Subject: Unmodified Message v2
	Date: Mon, 23 Dez 2314 12:23:42 +0200
	Message-Id: <RADNE23UOY.2KBV29410ATNY@example.org>

	This message is already in the archive.
EOF

"${ARCHIVE_MAIL}" current1→archive1 current2→archive2
check_maildir archive1 expected1
check_maildir archive2 expected2
