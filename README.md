<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Pubber Subber](#pubber-subber)
  - [Configuration](#configuration)
    - [Projects](#projects)
    - [Service Account](#service-account)
  - [Local Dev](#local-dev)
    - [Environment](#environment)
    - [Run](#run)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Pubber Subber

Publish and subscribe to PubSub messages over HTTP.

## Configuration

### Projects

Create the desired PubSub topics and subscriptions in GCP. Enter them into `config.json` in the project's root directory. That file stores a JSON object with the following keys:

```shell
# config.json
{
  "projects": [
    {
        "name":  "my-project"
      , "topics": ["topic1", "topic2"]
      , "subscriptions": ["test1, test2"]
    }
  ]
}
```

### Service Account

Copy GCP service account credentials into `secret/creds.json`.

## Local Dev

### Environment

Three environment variables are required:
Variable | Default
--------|---------
`GOOGLE_APPLICATION_CREDENTIALS` | `$PWD/secret/creds.json`
`CONFIG_PATH` | `$PWD/config.json`
`RECEIVE_ENDPOINT` | `http://host.docker.internal:8081`

### Run

1.  Run pubber-subber locally:

```bash
# Run locally
$ make start
# Or via docker
$ make docker
# Or build image locally and run
$ make docker-local
```

2. In another terminal, listen for HTTP requests on `RECEIVE_ENDPOINT` (defaults to localhost:8081):

```bash
$ node <<< '
let http = require("http"); let port = 8081;
console.log("listening on", port)
http.createServer((req, res) => {
  console.log("Received a message: " + req.method + " " + req.url)
  let str = ""
  req.on("data", chunk => str+=chunk)
  req.on("end", () => { console.log(str); res.end("done") })
}).listen(port)
'
```

3. Publish a PubSub message to pubber-subber, and watch it get passed to your server:

```bash
$  curl -X POST http://localhost:8080/publish/my-project/topic1 -d '
{
  "greeting": "hello from pubsub!"
}
'
```
