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
  - [x] Add timework
  - [x] Add file
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
# template base on: /api/new-task/new-requirement
bcli task create -t "[BE] Drop prod database" -T bug.json

# Edit task title, content
bcli task update 1234

# Edit task without open editor
bcli task update-headless 1234 --title='implement confirm popup' --content='update later'

# Add time work
# Add 2 hour to 11/05/2024 for task 23565
# template base on: /api/task-details/add-actual-effort-point
bcli task add-time 23565 -H 2 -T ~/.bcli/templates/add-time-dev.json -d 20241105

# Add effort point
# Add point to task 23565 with volume 2
# template base on: /api/save-req-job-detail
bcli task add-point 23565 -p 2 -T /home/tieu/.bcli/templates/add-point-dev.json

# Add file
bcli task add-file PRQ20241107000000188 -f ~/Downloads/image1.png

# More
bcli task -h

# Punch
bcli punch

# Generate shell completion. bash, zsh, fish, powershell
bcli completion bash
```
