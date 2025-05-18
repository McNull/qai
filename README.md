# qai

Simple AI query tool for the terminal using ollama or github copilot. Answers are short and to the point.

```bash
$ qai how can I scan 192.168.6.1 for the ports 22, 80, 8080
nmap -p 22,80,8080 192.168.6.1
```

## Installation

```bash
$ go install github.com/mcnull/qai@latest
```

Or download the latest release from the [releases page](https://github.com/McNull/qai/releases).

## Usage

```bash
Usage of qai:
  -color
        Enable colored output (default true)
  -config string
        Path to the config file (default "/home/null/.config/qai/config.json")
  -create-config
        Create a new config file with default values
  -debug
        Enable debug mode
  -debug-stream
        Enable debug response stream
  -github-login
        Create a new GitHub auth token
  -profile string
        Profile name
  -system string
        System prompt (default "The user is running a terminal in the following environment: {{.Platform}}.\nYour responses are {{.Verbose}}.")
  -verbose
        Enable verbose output
  -version
        Show version information
```

## Providers
Currently supports `ollama` and `github` providers. The behavior of the providers can be configured in the config file.

## Config
Default configuration file is `~/.config/qai/config.json`. 

```json
{
  "profile": "default",
  "system": "The user is running a terminal in the following environment: {{.Platform}}.\nYour responses are {{.Verbose}}.",
  "providers": {
    "ollama": {
      "model": "llama3.2",
      "url": "http://127.0.0.1:11434"
    },
    "github": {
      "model": "gpt-3.5-turbo",
      "token": ""
    }
  },
  "profiles": {
    "default": {
      "provider": "ollama"
    }
  }
}
```

### Profiles
Profiles are used to switch between different providers and configurations. You can create multiple profiles in the config file and switch between them using the `-profile` flag or change the default profile in the config file.

```json
{
  "profile": "my-custom-profile",
  // ...
  "profiles": {
    "my-custom-profile": {
      "provider": "github",
      "settings": {
        "model": "gpt-4",
        "system": "You're a toilet assistant. Answer the user's questions about toilets.",
      }
    }
  }
}
```

