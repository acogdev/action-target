# action-target monitoring service
## Description
This program will use TCP connections to monitor the status of provided services.

The status of the services will be displayed in an auto updating webpage at <http://localhost:8080>

## How to use
Once the project is cloned to the local machine running `make` from in the project directory will run the tests and build the binary. 

To run just the tests run `make test`.

To run just the build run `make build`.

To install as a systemd service run `make install`.

Once the binary has been built the program can be run like this:
```sh
./action-target --hosts "google.com,apple.com" -p 80 -i 5
```
This command will have the program monitor the google.com and apple.com on port 80 and will check if they are up every 5 seconds.
The webserver will also be started and can be accessed from <http://localhost:8080/>.

A config file can also be passed in using the `--config` cli parameter. This should be the fully qualified path to the config file. The config file should be in `TOML` format and look like the following:
```toml
hosts = ["google.com"]
port = 80
interval = 5
```
And the call should look like:
```sh
./action-target --config path/to/config.toml
```


## Decisions

### Language
Decided to use `go` because I wanted to learn the language and because it appears to be the language favored by the team. `go` also works well in this space of bridging systems programming and web programming. 

### Frameworks
I decided to use [cobra](https://github.com/spf13/cobra) for the commandline tooling to make creating a useful cli easier and to have the help, info, etc. built in.

I decided to use the [datatables](https://datatables.net) javascript project for much the same reason. It simplified implementation and made it easy to make it look good. It also provided a straight forward way to update the values periodically.

### Structure / Layout
I followed the layout of several other projects that use `cobra`, partly because I'm new to both `go` and `cobra` and because it seemed to be a good way of laying the files out.

