# biunzip

WIP

## Development Notes

- Support single-file mode with path (file path) and password flags.
- Support multiple-file mode with path (directory path) and CSV flags.
- ~~Check for missing and additional file when using multiple-file mode and provide interactive warnings.~~
- Validate zip files in the csv file.
- Accept the "Zip File" column or the first column as the filename in the CSV file.
- Accept the "Zip Password" column or the last column as the password in the CSV file.
- Set an empty password if the password flag is not provided.
- Extract each zip file to its own directory with the same name.
- Make integrations tests for both multiple and single commands with using Github workflows.

## How To Unzip All Of The Zip Files In A Directory

```bash
./biunzip --dir <directory_path> --csv <csv_file_path>
```

or

```bash
./biunzip -d <directory_path> -c <csv_file_path>
```

## How To Unzip A Single Zip File

```bash
./biunzip --file <zip_file_path>
```

or

```bash
./biunzip -f <zip_file_path>
```

## How To Unzip A Single Encrypted Zip File

```bash
./biunzip --file <zip_file_path> --password <zip_password>
```

or

```bash
./biunzip -f <zip_file_path> -p <zip_password>
```

## Help

```bash
./biunzip --help
```

or

```bash
./biunzip -h
```


