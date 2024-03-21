# windsor

A little cli tool to help you book tennis courts @ windsor tennis club.

## build

```bash
go build .
```

## running

Help:

```bash
./windsor --help
Usage: Windsor booking bot [flags]

A simple bot to book windsor tennis courts

Flags:
  -h, --help                                                    Show context-sensitive help.
      --username="thomas honey"                                 username
      --password="password"                                     password
      --participants=joe blogs...                               participants
      --day=STRING                                              specify date
      --rooms=43,44,45,46,47,...                                rooms
      --hour=20                                                 hour
      --area="24"                                               area
```

NB by default it tries to book 8 days later. So day is not needed.

Example usage:

```bash
./windsor --password=super_secret_password --username=joe blogs
```

## Bonus feature

There is also a python version of the go code as well.
