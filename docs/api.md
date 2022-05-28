# API Manual

Base command manual

```
NAME:
   musync - Sync spotify playlists

USAGE:
   musync [global options] command [command options] [arguments...]

COMMANDS:
   login    login into your spotify account
   whoami   show the current logged in account
   logout   logout from your spotify account
   sync, s  start a sync task between playlists / library
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help (default: false)
   --print-version, -V  print only the version (default: false)
```

### `login`

`musync login` manual

```
NAME:
   musync login - login into your spotify account

USAGE:
   musync login [command options] [arguments...]

OPTIONS:
   --clientid value        The client id from the spotify developer console, flag takes priority over MUSYNC_CLIENT_ID
   --clientsecret value    the client secret from the spotify developer console, flag takes priority over MUSYNC_CLIENT_SECRET
   --port value, -p value  port to run the server on, please make sure you match this on the spotify portal, flag takes priority over MUSYNC_CLIENT_ID (default: "8080")
```

### `logout`

`musync logout` manual

```
NAME:
   musync logout - logout from your spotify account

USAGE:
   musync logout [command options] [arguments...]

OPTIONS:
   --help, -h  show help (default: false)
```

### `whoami`

`musync whoami` manual

```
NAME:
   musync whoami - show the current logged in account

USAGE:
   musync whoami [command options] [arguments...]

OPTIONS:
   --help, -h  show help (default: false)
```

### `sync`

```
NAME:
   musync sync - start a sync task between playlists / library

USAGE:
   musync sync [command options] [arguments...]

OPTIONS:
   --dest value, -d value    Destination playlists, playlist to transfer everything into
   --quiet, -q               Will print nothing to stdout other than errors (default: false)
   --source value, -s value  Source playlist, type 'lib' to use your music library as the source (default: "lib")
```
