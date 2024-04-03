# Requirements

### Common

- Go >= 1.21.1

- VSCode or VSCodium with the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)

### Client

- Python >= 3.10

- GNU Make >= 3.0
  
  - Windows users can obtain a binary using `download_make7.bat` for Windows 7 and `download_make10.bat` for Windows 10 and 11

### Server

- [Fresh](https://github.com/gravityblast/fresh)
  
  - Install via: `go install github.com/pilu/fresh@latest`. Ensure `~/go/bin` is in your `$PATH` environment variable

- [go-enum](https://github.com/abice/go-enum)
  
  - Install via: `go install github.com/abice/go-enum@latest`. Ensure `~/go/bin` is in your `$PATH` environment variable

- A MongoDB server for storing data

- A Redis server for caching

- An email server for sending emails to users. See [this documentation page](./documentation/Installing%20A%20Dev%20Email%20Server.md) for a tutorial on setting one up for a development environment

### Optional

- MongoDB Compass: allows viewing the contents of the MDB database in a friendly GUI. It can be downloaded from [this link](https://www.mongodb.com/try/download/compasshttps://www.mongodb.com/try/download/compass).
- RedisInsight: allows viewing the contents of the Redis database in a friendly GUI. It can be downloaded from [this link](https://redis.com/redis-enterprise/redis-insight/).

# Running The Code

(TODO)
