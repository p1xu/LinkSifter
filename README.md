# LinkSifter

LinkSifter is a command-line tool that allows you to sift through a list of URLs, checking them against patterns to find matches.
Ideally, LinkSifter works perfectly with huge URL results that come from a web crawler or archive.


## Features

- Check URLs against patterns to find matches
- Support for various matching options (equal, contain, regex)
- Customize matching behavior (case-insensitive, URL components)
- Concurrent processing for faster execution


## Installation

### Prebuilt Binary
Download a prebuilt binary from the [releases page](https://github.com/p1xu/LinkSifter/releases). Choose the appropriate binary for your operating system and architecture.

### Using Go Compiler
If you have Go compiler installed, you can install LinkSifter using the following command:
```
go install github.com/p1xu/LinkSifter@latest
```

### Build from Source
Clone the repository:
```
git clone https://github.com/p1xu/LinkSifter.git
cd LinkSifter
```
Install the dependencies and build the application:
```
go get && go build .
```

## Usage

./LinkSifter [Flags]
```
Flags:
  -L, --all2lowercase     Convert URLs and patterns to lowercase
  -e, --equal             Check if patterns from the file are equal to the compared value (useful when comparing against a part of the URL)
      --filename          Check the Filename from the URL path only
  -h, --help              help for LinkSifter
  -i, --input string      File containing a list of URLs
  -l, --lowercase         Convert all URLs to lowercase
  -o, --output string     File to save results to
      --path              Check URL Path only (includes the filename)
      --rawpath           Check URL Path only with out decoding (includes the filename)
      --rawquery          Check the URL Querys only
  -r, --regex             Check using regex from the pattern file
  -t, --threads int       Number of concurrent threads (default 50)
  -v, --verbose           Enable verbose mode
  -w, --wordlist string   File containing a list of patterns
```

- Use the `-i` or `--input` flag to specify the input file containing the list of URLs.
- Use the `-o` or `--output` flag to specify the output file to save the results.
- Use the `-w` or `--wordlist` flag to specify the patterns file to check against.
- Customize the matching behavior and options using the available flags.

### Comparing Options
- Use the `-r` or `--regex` flag to convert the patterns list to regex list. If the `--verbose` flag is used, it will perform a precheck to ensure the regex is valid before using it.
- Use the `-e` or `--equal` flag for exact string comparison (useful when comparing against a specific component of the URL like `--path` flag).
- Use the `-l` or `--lowercase` flag to convert all URLs to lowercase.
- Use the `-L` or `--all2lowercase` flag to convert all URLs and patterns to lowercase.

### Scoop Options
- Use the `--path` flag to limit the search to the URL path only, for example: `/folder/υποφάκελο/file`.
- Use the `--rawpath` flag to limit the search to the URL path without decoding it, for example: `/folder/%CF%85%CF%80%CE%BF%CF%86%CE%AC%CE%BA%CE%B5%CE%BB%CE%BF/file`.
- Use the `--filename` flag to limit the search to the URL filename only, for example: `file`.
- Use the `--rawquery` flag to limit the search to the URL queries only, for example: `page=main&id=09283&callback=https%3A%2F%2F..`.

### Other Options
- Use the `-t` or `--threads` flag to specify the number of threads working simultaneously (default: 50).
- Use the `-v` or `--verbose` flag to enable more visualization during the sifting process.

**Note:** 
- If the ‘Scoop Options’ are not set, LinkSifter will compare the patterns against the full URL.
- If the `--regex` and `--equal` flag is not present, LinkSifter will check if the pattern is part of the URL using substring matching (strings.Contains).

### Example
```
LinkSifter -i urls.txt -o harvest.txt -w leaky-paths.txt --path
```

**Note:** I have edited the original list of leaky paths from https://github.com/ayoubfathi/leaky-paths

## Contribution

Contributions are welcome! If you find a bug or have suggestions for improvements, please open an issue or submit a pull request. Alternatively, you can reach out to me through my [Twitter](https://twitter.com/0_Pixu) account


## License

LinkSifter is released under the [MIT License](https://opensource.org/licenses/MIT). See the [LICENSE](LICENSE) file for more details.
