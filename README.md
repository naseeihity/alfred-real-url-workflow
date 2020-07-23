# alfred-real-url-workflow
An Alfred workflow for real-url, play your live playlist using IINA.

This is a project to learn Golang.

## Usage
### run from Command Line
```
# open and update url
# just run the code directly

# open by room id
play {platform} {room_id}

Example:
zb douyu 9999

# open room from local list, if not set platform, then will open all rooms
play {platform}

# add room to local list
add {platform} {room_id}

# open directly without update url
play
```

## Features
- Play all kinds of live in one window
- Display room info
- Cache real url in local file
- Read and Edit room list from Json file
- Concurrency get real url(Very fast even have many rooms)

## Support Platforms
- [x] bilibili
- [x] zhanqi
- [x] douyu
- [x] youku
- [ ] huya

## Reference
[wbt5/real-url](https://github.com/wbt5/real-url)