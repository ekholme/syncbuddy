# Syncbuddy

Syncbuddy is a small CLI tool to help synchronize files and directories.

## Installation

You can install syncbuddy with:

```bash
go install github.com/ekholme/syncbuddy@latest
```

## Usage

Sync files from a source directory to a destination directory.

```bash
syncbuddy sync --source /path/to/source --destination /path/to/destination
```

You can also use the shorter flags:

```bash
syncbuddy sync -s /path/to/source -d /path/to/destination
```

## License
[MIT](LICENSE)