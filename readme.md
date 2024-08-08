# Tailor

Multi-server, multi-file, remote tail like a boss in golang.

- 📟 Supports multiple servers and files.
- 🔑 Sources your SSH config.
- 📖 Supports tailing logs.
- 📋 Supports running commands.

## Usage

| Command                       | Description                        |
| ----------------------------- | ---------------------------------- |
| [`tailor tail`](#tailor-tail) | Tail logs from multiple servers.   |
| [`tailor run`](#tailor-run)   | Run a command on multiple servers. |

### `tailor tail`

Tail logs from multiple servers.

| Flag                | Description                                     |
| ------------------- | ----------------------------------------------- |
| `--addresses`, `-a` | The addresses of the servers to tail logs from. |
| `--files`, `-f`     | The files to tail logs from.                    |

```bash
tailor tail --addresses server1,10.0.0.1 --files /var/log/nginx/access.log /var/log/nginx/error.log
```

### `tailor run`

Run a command on multiple servers.
