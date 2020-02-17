# archive-mail

Maintains maildir archives synced with current maildirs.

## Motivation

I only store the last `N` messages in the maildir on my mail server, I
maintain a local archive which contains all mails I ever received. This
tool helps propagating new and modified messages from the maildir on my
server to my archive.

## Usage

Sample usage for archiving the `INBOX` and `GitHub` maildir:

	$ archive-mail mail/INBOX→/srv/nfs/archive/mail/INBOX \
		mail/GitHub→/srv/nfs/archive/mail/GitHub

This will propagate the following changes to the archive:

1. New messages from the current maildir, which were previously
   not tracked in the archive.
2. Changed flags, or file names in general, of messages already
   tracked in the maildir archive.
3. Location changes of messages in the archive. For example, messages
   moved between `new/` and `cur/` and messages moved between different
   maildirs. For example, between `INBOX` and `GitHub` in the example
   above.

The current maildir will never be modified. Messages deleted from the
current maildir will also not be deleted from the archive.

## Installation

After cloning this repository compile the software as follows:

	$ go build

Afterwards copy the binary to your `$PATH` or use `go install`.

## Tests

A minimal test suite is provided it can be invoked as follows:

	$ ./tests/run_tests.sh

## License

This program is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published by the
Free Software Foundation, either version 3 of the License, or (at your
option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
Public License for more details.

You should have received a copy of the GNU General Public License along
with this program. If not, see <http://www.gnu.org/licenses/>.
