# Pubber Subber

## Use

Create the desired GCP PubSub topics and subscriptions and enter them into `config.json`.

The config is a JSON object with `projects` referring to an array of GCP projects. Example:

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

Copy GCP service account credentials into `secret/creds.json`.

Run it:

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
