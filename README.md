<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Pubber Subber](#pubber-subber)
  - [Configuration](#configuration)
    - [Projects](#projects)
    - [Service Account](#service-account)
  - [Local Dev](#local-dev)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Pubber Subber

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

```bash
# Run pubber-subber in one terminal
$ make run
# In another terminal, run a server to act as a subscriber
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
)
# Publish a PubSub message and watch it get passed to your server
$  curl -X POST http://localhost:8080/publish/my-project/topic1 -d '
{
  "greeting": "hello from pubsub!"
}
'
```
