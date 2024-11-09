# BluePrint but in cli

## Todo
- [x] List task
- [x] View task
- [x] Auth
- [x] Punch
- [x] Create new task with template
- [x] Update task
  - [x] Update content
  - [x] Update title
  - [x] Add EP
  - [ ] Add timework
- [ ] Comment to task
- [ ] Improve
  - [ ] Default config (credential, cookie file path, ...)

## Install

Download pre-built packages for Linux, Windows and macOS are found on the
[Releases](https://github.com/letieu/bcli/releases/) page.

**Or** Install with `go install`
```bash
go install github.com/letieu/bcli@0.1.9
```

## Usage
```bash
# Help
bcli -h

# Login with credential file (contain user, password)
bcli auth login -f credential.txt
# Login with username, password directly
bcli auth login -u <username> -p <password>

# List task
bcli task list
# List task but cool
bcli task list -m
# Create new task with template bug.json
bcli task create -t "[BE] Drop prod database" -T bug.json
# Edit task title, content
bcli task update 1234
# More
bcli task -h

# Punch
bcli punch

# Generate shell completion. bash, zsh, fish, powershell
bcli completion bash
```
