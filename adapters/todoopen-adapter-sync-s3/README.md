# todoopen-adapter-sync-s3

Sync adapter that reads/writes `tasks.jsonl` to an S3-compatible object store.

Works with **AWS S3**, **Cloudflare R2**, **Backblaze B2**, **MinIO**, and any
other S3-compatible provider.

Implements the `todoopen.plugin.v1` protocol over stdin/stdout.

---

## Install

```sh
# npm
npm install -g @justestif/todoopen-adapter-sync-s3

# mise
mise use -g go:github.com/justEstif/todo-open/adapters/todoopen-adapter-sync-s3@latest
mise reshim

# build from source
git clone https://github.com/justEstif/todo-open.git
cd todo-open
go build -o todoopen-adapter-sync-s3 ./adapters/todoopen-adapter-sync-s3
```

---

## Configuration

Add the adapter to your `.todoopen/config.toml`:

```toml
[adapters.s3]
  bin  = "todoopen-adapter-sync-s3"
  kind = "sync"

[adapters.s3.config]
  bucket     = "my-tasks"
  key        = "tasks.jsonl"   # optional — default: tasks.jsonl
  region     = "us-east-1"     # optional — default: us-east-1
  endpoint   = ""              # optional — omit for AWS; set for R2/B2/MinIO
  access_key = "${S3_ACCESS_KEY}"
  secret_key = "${S3_SECRET_KEY}"
```

All `config` values support `${VAR}` environment variable expansion.
**Never store raw credentials in the config file.**

### Provider examples

**Cloudflare R2**
```toml
[adapters.s3.config]
  bucket     = "my-tasks"
  endpoint   = "https://${CF_ACCOUNT_ID}.r2.cloudflarestorage.com"
  region     = "auto"
  access_key = "${R2_ACCESS_KEY_ID}"
  secret_key = "${R2_SECRET_ACCESS_KEY}"
```

**MinIO (local)**
```toml
[adapters.s3.config]
  bucket     = "tasks"
  endpoint   = "http://localhost:9000"
  region     = "us-east-1"
  access_key = "${MINIO_ACCESS_KEY}"
  secret_key = "${MINIO_SECRET_KEY}"
```

**AWS S3 (default credential chain)**
```toml
[adapters.s3.config]
  bucket = "my-tasks"
  region = "us-east-1"
  # access_key/secret_key omitted — uses ~/.aws/credentials, IAM role, etc.
```

---

## Capabilities

| Method   | Description                                                   |
|----------|---------------------------------------------------------------|
| `push`   | Upload local `tasks.jsonl` to the configured bucket/key.      |
| `pull`   | Download object to local `tasks.jsonl`, overwriting it.       |
| `status` | Report local vs. remote file existence, size, and timestamps. |

### First-push behaviour

If `tasks.jsonl` does not exist locally, `push` returns a no-op message.
If the remote object does not exist yet, `pull` returns a no-op message —
no error is raised on first push.

---

## Building

```bash
go build ./adapters/todoopen-adapter-sync-s3
```

Place the resulting binary somewhere on your `PATH`.

---

## Protocol

The adapter speaks `todoopen.plugin.v1`: one JSON object per line on stdout
(handshake first, then responses); one JSON request per line on stdin.
