# ChromeDebugLnk

## Introduction

**ChromeDebugLnk** is a Go-based utility designed to modify Chrome browser shortcuts on the Windows desktop, taskbar, or user-specified paths to enable remote debugging mode. It operates with elevated User Account Control (UAC) privileges to bypass restrictions imposed by security software like 360, allowing seamless modification of shortcut properties. The tool can also restrict Chrome's incognito mode by modifying registry settings.

### Key Features
- **Modify Browser Shortcuts**: Adds remote debugging parameters (`--remote-debugging-port` and `--remote-allow-origins=*`) to Chrome, Edge, or Opera shortcuts.
- **Bypass Incognito Mode**: Optionally restricts Chrome's incognito mode via registry changes.
- **Customizable Paths**: Supports user-defined shortcut paths and usernames for targeted modifications.
- **Self-Deletion**: Optionally deletes the executable after execution for stealth.
- **UAC Check**: Ensures the program runs with administrative privileges.

### Usage
Run the program with the following command-line flags:
- `-a`: Password (must be `fuck360` to proceed).
- `-p`: Debugging port (default: `9222`).
- `-l`: Specify a custom shortcut name (e.g., `Google Chrome`).
- `-u`: Specify a username to target their desktop.
- `-path`: Specify a custom shortcut path.
- `-bypass`: Enable incognito mode restriction for Chrome.
- `-nobypass`: Remove incognito mode restriction.

**Example**:
```bash
ChromeDebugLnk.exe -a fuck360 -p 9222 -l "Google Chrome"
```