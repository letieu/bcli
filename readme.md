# BluePrint but in cli

## Todo
- [x] List task
- [x] View task
- [x] Auth
- [x] Punch
- [ ] Create new task with template
- [ ] Update task (Think about edit style)
  - [x] Update content
  - [x] Update title
  - [ ] Add EP, Timework to task
- [ ] Comment to task
- [ ] Improve
  - [ ] Default config (credential, cookie file path, ...)

## Usage
```bash
# Help
./bcli -h
# Login
./bcli auth login -f ~/.bcli/.credential.txt
# List task
./bcli task list
# List task but cool
./bcli task list -m
# Punch
./bcli punch
# Update task, view task detail, etc, use help for more detail
./bcli task -h
```

