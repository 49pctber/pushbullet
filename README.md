# Pushbullet

A CLI API wrapper for [Pushbullet](https://www.pushbullet.com/) notifications.

## Usage

1. Make a [Pushbullet](https://www.pushbullet.com/) account
2. Download the Pushbullet app to your mobile device
3. Install Go
4. Install the CLI executable using the command `go install github.com/49pctber/pushbullet@latest`
5. Set your `PUSHBULLET_TOKEN` environment variable to your Pushbullet API key.
6. Restart your terminal if necessary
7. Send your notification `pushbullet --title="[title]" --body="[body]"`
