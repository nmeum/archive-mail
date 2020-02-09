# archive-mail

Maintains a maildir archive synced with a current maildir.

## Motivation

I only store the last `n` messages in the maildir on my mail server, I
maintain a local archive which contains all mails I ever received. This
tools helps propagating new and modified messages from the maildir on my
server to my archive.

## Usage

Sample usage for archiving the `INBOX` and `GitHub` maildir:

	$ archive-mail mail/INBOX→/srv/nfs/archive/mail/INBOX \
		mail/GitHub→/srv/nfs/archive/mail/GitHub

The current maildir will never be modified.

## Tests

A minimal test suite is provided it can be invoked as follows:

	$ cd tests && ./run_tests.sh

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
