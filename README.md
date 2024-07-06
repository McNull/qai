# qai

Simple AI query tool for the terminal using ollama. Answers are short and to the point.

```bash
$ qai how can I scan 192.168.6.1 for the ports 22, 80, 8080
nmap -p 22,80,8080 192.168.6.1
```

## Usage

```bash
Usage of ./qai:
  -m, --model string    Model to use for generation (default "llama3")
  -o, --os string       Operating system (default "{your operating system}")
  -s, --system string   System message to use for generation
  -u, --url string      URL of the ollama service (default "http://localhost:11434")
  -v, --version         Print the version and exit
```

## Config
Default configuration file is `~/.config/qai/config.json`. 

```json
{
    "model": "llama3",
    "url": "http://localhost:11434",
    "system": "Your answers are always formal, short and to the point.\n\t\tYour answers never contain explanations or examples unless explicitly asked for."
}
```
