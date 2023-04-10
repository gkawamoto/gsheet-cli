# gsheet-cli

Fetch Google Sheets data from your command line.

# Installation

1. Install the `gsheet` command: `go install github.com/gkawamoto/gsheet-cli/cmd/gsheet@latest`
2. Create **OAuth client ID** Desktop app credentials on your Google Cloud Platform: https://console.cloud.google.com/apis/credentials
3. Download the JSON file and save it on `~/.config/gsheet/credentials.json`
4. Run `gsheet auth` and follow the steps
5. ???
6. Profit!

# Usage

```bash
$ gsheet help
Usage:
  gsheet [command]

Available Commands:
  auth        Checks whether your cli is authenticated correctly
  completion  Generate the autocompletion script for the specified shell
  get         Get data from a spreadsheet
  help        Help about any command

Flags:
  -d, --config-dir string         config directory (default "~/.config/gsheet")
  -c, --credentials-file string   credentials file location (default "~/.config/gsheet/credentials.json")
  -h, --help                      help for gsheet
  -r, --reauthenticate            reauthenticate

$ gsheet get --help
Get data from a spreadsheet

Usage:
  gsheet get [flags]

Flags:
  -h, --help                    help for get
  -s, --spreadsheet-id string   spreadsheet ID

Global Flags:
  -d, --config-dir string         config directory (default "~/.config/gsheet")
  -c, --credentials-file string   credentials file location (default "~/.config/gsheet/credentials.json")
  -r, --reauthenticate            reauthenticate
```

# Fetching data

```bash
$ gsheet get -s TW9uIEFwciAx_MCAwMzowMzoyMiBDRVNUIDIwMjMKAwM 'Sheet1!A1:C' | jq
[
  [
    "ID",
    "Name",
    "Age"
  ],
  [
    "1",
    "Gustavo",
    "33"
  ],
  [
    "2",
    "Ellen",
    "29"
  ]
]
```

# Contributing

Issues and PRs are welcome.

# License

MIT License