# biunzip

## DEPRECATED - PASSWORD PROTECTION IS REMOVED FROM OFF-NETWORK ZIP OUTPUTS 

[![Build](https://github.com/binalyze/biunzip/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/binalyze/biunzip/actions/workflows/build.yml)
[![Test](https://github.com/binalyze/biunzip/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/binalyze/biunzip/actions/workflows/test.yml)

biunzip is a command-line tool specifically designed to extract zip files generated by Binalyze Agent Off-Network. 

Please note that this tool is not intended for general-purpose zip file processing.

# Installation

You can download the latest release of biunzip from the [releases](https://github.com/binalyze/biunzip/releases) section. Alternatively, if you are familiar with the Go programming language, you can install it from the source by running the following command in your terminal.

```bash
go install https://github.com/binalyze/biunzip@latest
```

# Usage

After installing biunzip, you can use it to unzip zip files by executing the following commands in your terminal. There are two modes available: you can either unzip a single zip file or unzip zip files in a directory using a CSV file.

## Unzip A Single Zip File

You can unzip a single zip file by using the --file flag. Additionally, you have the option to specify a password with the --password flag if the zip file is encrypted.

### Unix

```bash
./biunzip --file zip_file_path --password zip_file_password
```

### Windows

#### cmd.exe

```shell
biunzip.exe --file zip_file_path --password zip_file_password
```

#### PowerShell

```powershell
.\biunzip.exe --file zip_file_path --password zip_file_password
```

## Unzip Zip Files In A Directory

You can unzip zip files in a directory with a CSV file by using --dir and --csv flags.

Please make sure the CSV file should include a header line with a minimum of two columns. The first column or the "File Name" labeled column should contain the names of the zip files. The last column or the "Zip Password" labeled column should contain the password if the zip files are encrypted.

**Sample CSV File:**
```csv
File Name,Zip Password
file_1.zip,password_1
file_2.zip,password_3
```

The sample CSV file has the required two columns. It is also possible for the CSV file to have extra columns, which will not cause an error for biunzip.

### Unix

```bash
./biunzip --dir dir_path --csv csv_file_path
```

### Windows

#### cmd.exe

```shell
biunzip.exe --dir dir_path --csv csv_file_path
```

#### PowerShell

```powershell
.\biunzip.exe --dir dir_path --csv csv_file_path
```

## Help

To view a detailed help message, run the following command in your terminal.

### Unix

```bash
./biunzip --help
```

### Windows

#### cmd.exe

```shell
biunzip.exe --help
```

#### PowerShell

```powershell
.\biunzip.exe --help
```

# License

biunzip is licensed under the [Apache License](LICENSE).
