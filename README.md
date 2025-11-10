# ws

CLI tool and Go library to interact with Wavin Sentio floor heating systems.

## Installation

### Homebrew (macOS and Linux)

```sh
brew install zmoog/ws/ws
```

### Go Install

```sh
go install github.com/zmoog/ws/v2@latest
```

## Usage

List devices in your Wavin Sentio account:

```sh
$ ws devices list

Name                          | Serial Number  | Type           | Firmware | Last Heartbeat
devices/abcdefghijklmnopqrstu | 98765432109876 | TYPE_SENTIO_CCU| 17.2.1   | 2025-01-15T14:32:18Z
```

List rooms in a device:

```sh
$ ws rooms list --device-name devices/abcdefghijklmnopqrstu

Name                  | Temperature state | Temperature (desired) | Temperature (current) | Humidity (current) | Dehumidification state
Living Room           | TEMPERATURE_STATE_IDLE | 22.0                | 21.5                  | 45.2               | DEHUMIDIFIER_STATE_IDLE
Kitchen               | TEMPERATURE_STATE_IDLE | 20.0                | 19.8                  | 52.1               | DEHUMIDIFIER_STATE_IDLE
Bedroom               | TEMPERATURE_STATE_IDLE | 18.5                | 18.9                  | 48.7               | DEHUMIDIFIER_STATE_IDLE
Bathroom              | TEMPERATURE_STATE_IDLE | 23.0                | 22.4                  | 58.3               | DEHUMIDIFIER_STATE_IDLE

Outdoor temperature: 15.2
```

## Configuration

### Authentication

The tool uses a configuration file to store your credentials and tokens. The file is located in `~/.ws/config`.

```sh
$ cat ~/.ws/config
username: john.doe@example.com
password: secretpassword123
# web_api_key: AIzaSyBlAtNI7-2jitPul9I-O4EZcT-n0sIay-g  # Optional: uses default if not specified
api_endpoint: https://blaze.wavinsentio.com/wavin.blaze.v1.BlazeDeviceService
output: table
```

You can also set the credentials using:

- Command line flags: `--username`, `--password`, `--web-api-key`, `--api-endpoint`
- Environment variables: `WS_USERNAME`, `WS_PASSWORD`, `WS_WEB_API_KEY`, `WS_API_ENDPOINT`

### Firebase Web API Key

The `web_api_key` identifies the Wavin Sentio app on Firebase and is not a secret. The tool comes with a default key that should work for most users. You only need to specify a custom key if:

- The default key stops working (rare)
- You're using a custom or enterprise version of the Wavin Sentio app

The default value is: `AIzaSyBlAtNI7-2jitPul9I-O4EZcT-n0sIay-g`

### Output formats

The tool supports several output formats:

- `table`: (default) prints the output in a table format
- `json`: prints the output in JSON format

You can change the output format using the `--output` flag.

```sh
$ ws devices list --output json
[
  {
    "name": "devices/abcdefghijklmnopqrstu",
    "createTime": "2024-03-15T10:45:22.123456Z",
    "updateTime": "2025-01-15T14:32:18.987654321Z",
    "serialNumber": "98765432109876",
    "registrationKey": "A1B2C-D3E4F-5G6H",
    "firmwareAvailable": "17.2.1",
    "firmwareInstalled": "17.2.1",
    "type": "TYPE_SENTIO_CCU",
    "lastHeartbeat": "2025-01-15T14:32:18.987654321Z"
  }
]
```

## Migration from v1

This is version 2 of the tool, which uses the new Wavin Sentio backend. If you're upgrading from v1:

### Breaking Changes

- Command changed: `ws locations list` → `ws devices list`
- Flag changed: `--location-id` → `--device-name` 
- Module path: `github.com/zmoog/ws/v2`
- New optional config: `web_api_key` (has sensible default)

### Migration Steps

1. Uninstall the old version: `go clean -i github.com/zmoog/ws`
2. Install the new version: `go install github.com/zmoog/ws/v2@latest`
3. Update any scripts to use `devices` instead of `locations` and `--device-name` instead of `--location-id`
4. The `web_api_key` is now optional and has a default value
