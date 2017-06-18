# Murmamau
Murmamau collects data from a system and stores it inside of a TAR file for further analyses.
It fetches:
- Environment settings
- List of all users (all UIDs)
- Host stats
- Disk stats
- Memory stats
- CPU stats
- Network stats
- Information on all running proccesses
- Information on all files on the file system
- It saves files
  - Currenlty open files (based on process information)
  - Path or filename on a list
  - filenames on each home directory
  - filenames matching a pattern list (glob)


# Running
Run the binary and give path as argument (if not given "/" is used)


# Installation
Just copy binary "murmamau" to zour target

# Building
Compile static binary for current plattform:
```
CGO_ENABLED=0 go build --ldflags '-extldflags "-lm -lstdc++ -lvdso -static"'
```

# Configuation
Edit YAML configuration in "config/config.go"
- colors: enable colors in output?
- debug: enabled debug output?
- everyxstatus: show every x files a status message
- maxsize: maximum file size for saving into TAR file
- searchlist: list to download
  - can be path
  - can be filename
  - can start with "~/" for checking it in every home directory
- globlist: list of glob search pattern
- ignorelist: list of file extensions for ignoring
- blacklist: list of filenames for ignoring (NOT to save into TAR file)
commands: list of commands to be executed per plattform


# Issues
- tar shows error message "tar: A lone zero block at": use "tar -i" option for ignoring zeros
- Murmamau needs many resources: use nice
- Turn on debugging inside of the YAML configuration often helps

