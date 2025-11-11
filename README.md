# GitEase: Terminal-Based Git Assistant

**GitEase** is a simple, terminal-based Git assistant built in Go using the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework. 

It was created as part of a Go learning journey and turns everyday Git commands — such as staging, committing, and pushing — into an **interactive terminal interface**.

---

##  Installation

### Linux / macOS

Use `curl` to download and execute the installation script:

```bash
curl -sSL https://raw.githubusercontent.com/dinesh00509/gitease/main/install.sh | bash
```

### Windows (PowerShell)
Open PowerShell and run the following command to download and execute the installation script:

 ```bash
iwr -useb https://raw.githubusercontent.com/dinesh00509/gitease/main/install.ps1 | iex
```
Once installed, verify the installation:
```bash
gitease --version
```

### Usage
After installation, simply run:
```bash
gitease --run
```

