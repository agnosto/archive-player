# Archive Player

<p align="center">
  <img src="./build/appicon.png" alt="App Icon" width="128" height="128">
</p>

<p align="center">
  A <a href="https://archive.ragtag.moe/player">ragtag</a>-like player <br> with built-in integration for 
  <a href="https://github.com/agnosto/fansly-scraper">fansly-scraper</a> livestream VODs.
</p>

<p align="center">
    <small>built with <a href="https://wails.io">wails</a> (go version of tauri/electron, i know, sorry)</small>
    <br>
    <small>I probably won't add any other integrations unless it's for something I end up using</small>
</p>

---

## Fansly Scraper Integration

Archive Player includes integration with the Fansly Scraper tool, allowing you to easily browse and watch recorded livestreams with their associated chat logs.

### Setup

1. Open the Settings panel and click "Configure" under Fansly Integration
2. Provide the path to your Fansly Scraper config file (usually at `~/.config/fansly-scraper/config.toml` on Mac/Linux or `%APPDATA%\fansly-scraper\config.toml` on Windows)
3. Provide the path to the folder containing the `downloads.db` file (usually the `save_location` specified in your Fansly Scraper config)
4. Click "Save" to complete the setup

### Usage

1. Click the "Fansly Streams" button in the top menu to open the Fansly browser
2. Browse available streams, filter by model name, or search for specific content
3. Click on a stream to load it along with its chat (if available)
4. Enjoy watching the stream with synchronized chat replay

### Features

- Automatically pairs videos with their chat logs and contact sheet thumbnails
- Displays stream information including model name, date, and duration
- Allows filtering and searching by model name or filename
- Clip creation tool to extract segments (up to 5 minutes) from videos (requires FFmpeg)
- Preserves all Archive Player features like theater mode and chat display options

## Development

### Requirements

- [FFmpeg](https://ffmpeg.org/download.html) is required for the clip creation feature

### Chat Format

The player supports chat files in a format similar to what [chat-downloader](https://pypi.org/project/chat-downloader/) provides:

```json
{
  "message_id": "xxxxxxxxxx",
  "message": "actual message goes here",
  "message_type": "text_message",
  "timestamp": 1613761152565924,
  "time_in_seconds": 1234.56,
  "time_text": "20:34",
  "author": {
    "id": "UCxxxxxxxxxxxxxxxxxxxxxxx",
    "name": "username_of_sender",
    "images": [...],
    "badges": [...]
  }
}
```

Or from the fansly-scraper format:

```json
{
  "message_id": "xxxxxxxxxx",
  "message": "message content here",
  "message_type": "text_message",
  "timestamp": 1746116588125,
  "time_in_seconds": 4.686060153,
  "time_text": "00:04",
  "author": {
    "id": "xxxxxxxxxxxxxxxxxxxxxxx",
    "name": "display name of sender",
    "is_creator": true,
    "tier_info": {}
  },
  "raw_data": "...",
  "received_at": "2025-05-01T09:23:08.307450813-07:00"
}
```

### Building from Source

1. Install [Go](https://golang.org/doc/install) (1.24 or later)
2. Install [Wails](https://wails.io/docs/gettingstarted/installation)
3. Install [FFmpeg](https://ffmpeg.org/download.html) for clip creation functionality
4. Clone the repository
5. Run `wails build` to build the application

### Creating New Integrations

To add a new integration:

1. Create a new package in the `backend/integrations` directory
2. Implement the required service methods (similar to the Fansly integration)
3. Add your integration to the `Manager` in `backend/integrations/manager.go`
4. Create a UI component for your integration
5. Add the necessary methods to the `App` struct to expose your integration to the frontend

## Support the Project

If you find this tool useful and would like to support its development, you can donate using the following crypto addresses:

<table>
  <tr>
    <td align="center"><strong>Bitcoin (BTC)</strong></td>
    <td align="center"><strong>Solana (SOL)</strong></td>
  </tr>
  <tr>
    <td align="center">
      <img src="./assets/btc_qr.png" alt="Bitcoin QR Code" width="200"/>
      <p><code>bc1q0e78wrtc9ezp6tqv000wfewgqf2ue4tpzdk7ee</code></p>
    </td>
    <td align="center">
      <img src="./assets/sol_qr.png" alt="Solana QR Code" width="200"/>
      <p><code>Bv3kYZcwSTHXAQtnPddTF27D3F6Gc29v2MfFLqmGF6Gf</code></p>
    </td>
  </tr>
</table>

