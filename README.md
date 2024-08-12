# Kleos Logging Package

Kleos is a simple logging with a few very opinionated ideas behind it:

* Logging is for humans to read.  If you want to measure output, use a metric.
* You shouldn't have to pass a logger into every function to log messages.
* Various log levels are overly complicated and require too much effort.
* Debug log messages are for developers; other logs are for devops.
* Logs should be fast, but don't prematurely optimize.
* Structured logging is the easiest to parse, index, and search.

For Kleos, here's how you log a message:

    kleos.Log("This is my message")

Kleos tries to apply structured information to logging in a way that feels natural, so
developers can worry about conveying information rather than how they're conveying the
information. For example, rather than worrying about debug, info, warn, or error, how
about simply:

    kleos.Error(err).Log("Tried to do X and it didn't work")

Because you passed in an `error`, Kleos automatically turns this into an error message for
you on output. Without the call to `Error`, the message comes out as an info message.

Or maybe you've got some info you want to share with developers, but not something you
want to log in production?

    kleos.V(1).Log("Didn't expect to reach this line of code")

By default, Kleos won't output any log messages with verbosity attached. So you won't see
this message in production unless you turn up the verbosity.

Kleos also encourages structured logging. Rather than use string subsitution, which can
not only be slow, but difficult to filter for if you're watching logs, Kleos prefers you
use a consistent, static string message coupled with key/value pair properties.

    kleos.With(kleos.Fields{
        "email": "bob@nowhere.com", 
        "age": 37
    }).Log("Registered new user")

## Global Logger

Because log output is headed for stdout or a file, from my perspective it doesn't make a
ton of sense to be passing a logging variable around to every single function. Instead,
Kleos leverages a global variable. Yes, a package-level global variable. So rather than
pass a `logger` variable around in your function, you can just write:

    kleos.Log("Here is my log message.")

The logger ensures that messages are streamed in the order they're received, so it's ok
to call Kleos from different threads.

Kleos also supports local loggers. Internally, it uses a `kleos.Logger` to manage the
global logger. If you like, you may leverage that to create your own loggers and pass
them directly into functions or store them in structs.

    logger := kleos.New(outputWriter)
    logger.Log("This is a log message.")

## Fields

Kleos encourages structured logging. It doesn't support `fmt.Printf`-style output.
Instead, use `kleos.Fields` to pass data, like a JSON object:

    kleos.With(kleos.Fields{
        "email": "bob@nowhere.com",
        "url": "https://website.com",
    }).Log("Here's something a user did with this web site."

## Output

Kleos outputs to `os.Stdout` by default. It's a simple text writer. The output looks
something like this:

![Log text output](docs/text_output.png?raw=true "Log text output")

To use this but perhaps pass the output to a different stream, such to a file:

    file, _ := os.Create("/tmp/app.log")
    kleos.SetOutput(kleos.NewTextOutput(file))

Kleos offers two other output types:  color and JSON. Color output is primarily meant for
development. It's like text output, but colorized the output. The same log messages
above, colorized, require:

    kleos.SetOutput(kleos.NewColorOutput(os.Stdout))

Which will look something like this:

![Log color output](docs/color_output.png?raw=true "Log color output")

JSON output is meant to be used in production.

    kleos.SetOutput(kleos.NewJSONOutput(file))

It looks like a JSON document:

![Log JSON output](docs/json_output.png?raw=true "Log JSON output")

There's also a preliminary [LogStash](https://www.elastic.co/logstash) writer. This will
send your log output to ElasticSearch. Note that this is a writer, so use this with
JSON output:

    logstash := kleos.NewLogstashWriter(host, 5*time.Second)
    kleos.SetOutput(kleos.NewJSONOutput(logstash)

A common pattern I use is to configure a "dev mode" on startup.  By default, a project
using Kleos starts in a "dev mode."  This outputs colorized log messages to `os.Stdout`.
In production, I enable an environment variable which outputs JSON objects to a log file,
`os.Stdout` (maybe in a Kube cluster), or to Logstash.

## Developer Logs

Some log messages only make sense for developers. Kleos handles these through verbosity.
Verbosity can be adjusted at startup or really at any time. For example, if you're
building a server, you can include an endpoint to adjust the verbosity level remotely.

Verbosity starts at 1 and can go up to any number (typically 4 is the max). I've used
the following as a guideline for verbosity levels:

1. General purpose messages like "User created" or "Rendered page."
2. A bit more detail about a process, such as the steps in a process: "requested record X",
   "updated X with Y", "cached X to the datastore."
3. Low level details, such as JSON objects or SQL queries.
4. Very specific details, such as large XML documents, or every single HTTP request body.

Obviously you can create your own scale, but I've found the above works pretty well for
most applications.

To set the verbosity level, simply:

    kleos.SetVerbosity(3)

This will tell Kleos output to include any log messages with verbosity 1, 2, or 3.

