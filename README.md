# ws

CLI tool and Go library to interact with Wavin Sentio floor heating systems.

## Installation

```sh
go install github.com/zmoog/ws
```

## Usage

List locations in your Wavin Sentio account:

```sh
$ ws locations list

Ulc   | Registration     | Serial Number | Mode  | Vacation On | Outdoor Temperature | DST
12345 | 11111-22222-3333 | 1122334455    | ready | false       | 4.0                 | true
```

List rooms in a location:

```sh
$ ws rooms list --location-id 12345

Name                  | Status  | Temperature (desired) | Temperature (current) | Humidity (current)
1 CUCINA              | idle    | 21.0                  | 21.2                  | 53.2
2 STUDIO              | idle    | 18.0                  | 21.7                  | 58.1
4 CAMERA RAGAZZI      | idle    | 19.0                  | 20.3                  | 63.8
3 BAGNO 1             | heating | 21.0                  | 20.6                  | 67.9
5 BAGNO 2             | heating | 20.0                  | 19.8                  | 63.5
6 CAMERA MATRIMONIALE | idle    | 18.5                  | 21.4                  | 58.7
```

## Configuration

### Authentication

The tool uses a configuration file to store your credentials and tokens. The file is located in `~/.ws/config`.

```sh
$ cat ~/.ws/config
username: myusername
password: mypassword
output: table
```

You can also set the username and password:

- using the `--username` and `--password` flags, or
- by setting the `WS_USERNAME` and `WS_PASSWORD` environment variables.

### Output formats

The tool supports several output formats:

- `table`: (default) prints the output in a table format
- `json`: prints the output in JSON format

You can change the output format using the `--output` flag.

```sh
$ ws locations list --output json
[
  {
    "ulc": "12345",
    "registrationKey": "11111-22222-3333",
    "serialNumber": "1122334455",
    "attributes": {
      "mode": "ready",
      "vacationOn": false,
      "vacationUntil": null,
      "outdoor": {
        "temperature": 4.2
      },
      "dst": true
    }
  }
]
```
