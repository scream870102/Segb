# Segb

This is a discord bot based on [discordgo](https://github.com/bwmarrin/discordgo).

---

## Main Features

- Manage uploaded GIFs by storing their URLs and tags in a google spreadsheet
  - Actually, other contents are also supported, whatever a google spreadsheet cell can save
- Post the corresponding GIF (content) when queried

## Prerequisites

- A google spreadsheet
- A google service account
  - With its json key
  - Have read+write permission to the above google spreadsheet
- A discord Bot
  - With its token
  - Have at least view+send message permissions to your discord server/channel
- C# package
  - [discordgo](https://github.com/bwmarrin/discordgo)
  - [uuid](https://github.com/github.com/google/uuid)
  - Google.Apis.Sheets.v4

- A config file name with `config.json` in root with following content

``` json
{
 "token": "discord-token",
 "prefix": "prefix",
 "delay": 15,
 "channel": "channel id",
 "guild": "guild id"
}
```

- A token file name with `token.json` in root with following content

``` json
{
 "Email": "mail",
 "PrivateKey": "key",
 "TokenURL": "url",
 "Scopes": [
  "scope"
 ]
}
```

## Usage

> The following commands should be called with a prefix `~`

- `help|h`
  - Show help message
    - `~help`
- `add|a` [TAG1 TAG2 ... TAGN ]  CONTENT
  - Add new content
    - `~add Oozono Momoko Kawaii`
- `list|l`
  - Show content added in this server
  - `~list` [TAG1 TAG2 ... TAGN ]
  - `~list` [Author|author]=[@User] [TAG1 TAG2 ... TAGN ]
    - `~list oozono momoko`
    - `~list Author=@scream870102 oozono momoko`
    - `~list Author=@scream870102`
- `show|s` [TAG1 TAG2 ... TAGN ]
  - Show the content filtered by TAGS in the database
  - At least one `TAG` should be provided
    - `~show Oozono Momoko`
