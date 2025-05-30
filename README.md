go-picam-stream
====

## Binary

### Installation

#### Latest
```sh
curl -fsSL https://github.com/ylallemant/go-picam-streamer/raw/main/install.sh | bash
```

#### Specific Version
```sh
curl -fsSL https://github.com/ylallemant/go-picam-streamer/raw/main/install.sh | bash -s -- --version="<version>"
```

### Upgrade

```sh
picam-streamer upgrade [--force]
```

## What could be the plan

- stream
- non-blocking
- start/stop/pause
- take a picture
- take a video
- crate timelapse
- timelaps
- RTMP


## Resources

- https://medium.com/wisemonks/implementing-websockets-in-golang-d3e8e219733b
- RTMP
  - https://github.com/bluenviron/mediamtx
  - https://github.com/AgustinSRG/rtmp-server
  - https://github.com/yutopp/go-rtmp
