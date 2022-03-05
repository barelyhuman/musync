<p align="center">
    <img height="320" src="/assets/terminal.gif"/>
</p>

<h1 align="center">
musync
</h1>
 
<p align="center">
Simple utility to sync your spotify library with a playlist. 
</p>

## Features

- Sync Queues
- Batched Processing (less heavy on the RAM)

## Why

There's enough reasons to sync your spotify library to a playlist, this could go from sharing the playlist to keeping the library small and clean while still having a playlist as an archive.

## Install

You can download from the [releases](/releases) page or using the command below

```sh
# golang not installed?
curl -sf https://goblin.reaper.im/github.com/barelyhuman/musync | sh

# golang installed and > 1.15
go install github.com/barelyhuman/musync
```

## Usage

This is a very developer centric tool and would need a few steps to get it to up and ready.

1. Get the `CLIENT ID` and `CLIENT SECRET` from a [developer.spotify.com](https://developer.spotify.com) account
2. Set a redirection url on the developer portal and make sure you remember the port you used. (eg: `localhost:8080`)
3. Create a file `musync.yml` in `~/.config/musync/`

```sh
touch ~/.config/musync/musync.yml
```

4. copy the template from `config.template.yml` in this repo and paste it in `~/.config/musync/musync.yml`

5. Fill in the values there

```yml
client_id: "CLIENT_ID"
client_secret: "CLIENT_SECRET"

# the id of the playlist to sync the library to, you can get this from the https://open.spotify.com url when you have the playlist open
playlist: ""

# the port you set on the developer.spotify.com website when creating a new app
port: 8080
```

6. Once done with everything above, you can just run `musync` from anywhere and it should guide you through a authentication process or just start syncing process

```sh
musync
```

## LICENSE

[MIT](LICENSE)
