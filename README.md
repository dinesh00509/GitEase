 ### GitEase

This project was created as part of my Go learning journey. GitEase is a simple, terminal-based Git assistant built in Go using the Bubbletea framework.
It turns everyday Git commands — like staging, committing, and pushing — into an interactive terminal interface.

### Installation
#### Option 1 — Install with Go
If you have Go installed, you can grab GitEase directly with:

go install github.com/dinesh00509/GitEase@latest

Once installed, run it:

gitease

#### Option 2 — Build from Source

Clone the repo and build it manually:

git clone https://github.com/dinesh00509/GitEase.git
cd GitEase
go build -o gitease
sudo mv gitease /usr/local/bin/


Now you can run:

gitease
