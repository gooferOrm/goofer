# Installing the Goofer ORM CLI

This guide will walk you through the process of installing the Goofer ORM CLI on your system.

## Prerequisites

Before installing the Goofer ORM CLI, make sure you have the following prerequisites:

- Go 1.21 or later installed on your system
- Git for cloning the repository (if installing from source)

## Installation Methods

There are several ways to install the Goofer ORM CLI:

1. [Using Go Install](#using-go-install) (recommended)
2. [Using Homebrew](#using-homebrew) (macOS and Linux)
3. [Downloading Binaries](#downloading-binaries)
4. [Building from Source](#building-from-source)

## Using Go Install

The easiest way to install the Goofer ORM CLI is using Go's built-in installation tool:

```bash
go install github.com/gooferOrm/goofer/cmd/goofer@latest
```

This will download, compile, and install the latest version of the CLI to your `$GOPATH/bin` directory. Make sure this directory is in your system's `PATH` to use the `goofer` command from anywhere.

To verify the installation, run:

```bash
goofer version
```

You should see output similar to:

```
Goofer ORM v0.1.0
```

## Using Homebrew

If you're using macOS or Linux with Homebrew, you can install the Goofer ORM CLI using:

```bash
brew tap gooferOrm/goofer
brew install goofer
```

To verify the installation, run:

```bash
goofer version
```

## Downloading Binaries

You can download pre-compiled binaries for your operating system from the [GitHub Releases](https://github.com/gooferOrm/goofer/releases) page.

1. Go to the [Releases](https://github.com/gooferOrm/goofer/releases) page
2. Download the appropriate binary for your operating system and architecture
3. Extract the archive
4. Move the `goofer` binary to a directory in your `PATH`

For example, on macOS or Linux:

```bash
# Download the binary (replace X.Y.Z with the version and OS/ARCH with your system)
curl -L https://github.com/gooferOrm/goofer/releases/download/vX.Y.Z/goofer_X.Y.Z_OS_ARCH.tar.gz -o goofer.tar.gz

# Extract the archive
tar -xzf goofer.tar.gz

# Move the binary to a directory in your PATH
sudo mv goofer /usr/local/bin/

# Verify the installation
goofer version
```

On Windows, you can download the ZIP file, extract it, and add the directory to your `PATH` environment variable.

## Building from Source

If you want to build the CLI from source, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/gooferOrm/goofer.git
cd goofer
```

2. Build the CLI:

```bash
go build -o goofer ./cmd/goofer
```

3. Move the binary to a directory in your `PATH`:

```bash
# On macOS/Linux
sudo mv goofer /usr/local/bin/

# On Windows, move it to a directory in your PATH
```

4. Verify the installation:

```bash
goofer version
```

## Shell Completion

The Goofer ORM CLI supports shell completion for Bash, Zsh, Fish, and PowerShell.

### Bash

Add the following to your `~/.bashrc` file:

```bash
source <(goofer completion bash)
```

### Zsh

Add the following to your `~/.zshrc` file:

```bash
source <(goofer completion zsh)
```

### Fish

Add the following to your `~/.config/fish/config.fish` file:

```fish
goofer completion fish | source
```

### PowerShell

Add the following to your PowerShell profile:

```powershell
goofer completion powershell | Out-String | Invoke-Expression
```

## Updating the CLI

To update the Goofer ORM CLI to the latest version:

### Using Go Install

```bash
go install github.com/gooferOrm/goofer/cmd/goofer@latest
```

### Using Homebrew

```bash
brew upgrade goofer
```

### Using Binaries or Source

Follow the same installation process with the new version.

## Troubleshooting

### Command Not Found

If you get a "command not found" error when trying to run `goofer`, make sure:

1. The installation was successful
2. The directory containing the `goofer` binary is in your `PATH`

You can check your `PATH` with:

```bash
echo $PATH
```

### Permission Denied

If you get a "permission denied" error, make sure the binary is executable:

```bash
chmod +x /path/to/goofer
```

### Other Issues

If you encounter any other issues, please check the [GitHub Issues](https://github.com/gooferOrm/goofer/issues) page or create a new issue.

## Next Steps

Now that you have installed the Goofer ORM CLI, you can:

- Learn about the [available commands](./commands)
- See how to [configure the CLI](./config)
- Start [creating migrations](./migration)
- Begin [generating code](./generate)