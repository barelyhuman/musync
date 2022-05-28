# musync

> CLI tool / package to create sync tasks for public spotify playlists and user libraries

## Highlights

- Fast and Cross Platform
- Sync Queues
- Concurrent Batches Processing
- Modular API to extend and use to create additional tooling on top

## Motivation

Cleaning up my spotify library is something I do often and then regret when I can't find names of tracks I used to listen to, so I maintain a capsule playlist which is kept in sync with the library as often as I can, this can be redundant and so it was easier to write a tool that would do that for me.
That's when musync v0.0.1 was made with the simple functionality to just run a sync job based on the provided config. The newer version now is a full fledged cli instead of a complicated go lang binary.

## Philosophy

The CLI serves to be of the unix philosophy, that is to do one and do it well. Though, the remaining of the package due to the structure is reusable in your projects if you plan to build something similar. No limitations there.

> **Note**: if you are using the package and have something that you wish the package provided, please raise an issue

## Usage

#### Install

The CLI is available as binaries on the [releases](/releases) page and if you can't find one for your arch, you can try using [goblin](https://goblin.reaper.im) to build a binary for you as shown below

```sh
curl -sf https://goblin.reaper.im/github.com/barelyhuman/musync | sh
```

If you have `golang>=1.8` then you can request the binary from `go get`

```sh
go get -u github.com/barelyhuman/musync
```

#### Next Steps

This is a very developer centric tool and would need a few steps to get the needed data, this is done so that you can use your own spotify apps to own your playlist id, if you don't mind the playlist id's being stored on someone else's server, we'll have a service for this up soon. Which you can use instead of the below setup.

**Account Setup**

1. Go to [developer.spotify.com](https://developer.spotify.com) and login with your account and create an app, this will give you the `CLIENT_ID` and `CLIENT_SECRET`, we'll need it in a few steps
2. While creating the app, you set a redirection url which will be used by the CLI to handle the verification of your account, you can set it as `localhost:8080/callback` as `8080` is the default port used by the `musync`, if you do change the port to something else like `3000`, make sure you pass the `-p` flag to the `login` command with the port value.

eg: `musync login -p 3000`

**CLI Usage**

```sh
# create a login request
musync login --clientid <client-id-from-spotify> --clientsecret <client-secret-from-spotify> -p 8080

# check if the login was successful
musync whoami

# sync between library and playlist
musync sync -s "lib" -d xxxxxxx

# sync between public playlist and playlist that you have edit access to
musync sync -s "lib" -s xxxxxxx -d xxxxxxx
```

You can find `xxxxxxx` (which is the playlist's identifier) from the spotify playlist urls `https://open.spotify.com/playlist/1fSExQf33fHoAtqnGY99tl` here `1fSExQf33fHoAtqnGY99tl` is what you'll pass to musync.

for more commands and flags check the [API documentation](/docs/api.md)

## LICENSE

[MIT](LICENSE)
