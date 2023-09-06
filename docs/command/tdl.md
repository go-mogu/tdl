## tdl

Telegram Downloader, but more than a downloader

### Options

```
      --debug                        enable debug mode
  -h, --help                         help for tdl
  -l, --limit int                    max number of concurrent tasks (default 2)
  -n, --ns string                    namespace for Telegram session
      --ntp string                   ntp server host, if not set, use system time
      --proxy string                 proxy address, only socks5 is supported, format: protocol://username:password@host:port
      --reconnect-timeout duration   Telegram client reconnection backoff timeout, infinite if set to 0 (default 2m0s)
  -s, --size int                     part size for transfer, max is 512*1024 (default 524288)
      --test string                  use test Telegram client, only for developer
  -t, --threads int                  max threads for transfer one item (default 4)
```

### SEE ALSO

* [tdl chat](tdl_chat.md)	 - A set of chat tools
* [tdl login](tdl_login.md)	 - Login to Telegram
* [tdl version](tdl_version.md)	 - Check the version info

