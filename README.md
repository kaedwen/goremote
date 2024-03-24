# goremote

This small project should enable you to perform tasks remotely on a machine.

### Config

Most properties can be specified as arguments and environment variables. Setting the argument `--config-file` or the environment variable `CONFIG_FILE` will also read the configuration from there with highest priority. Default location `($XDG_CONFIG_HOME|$HOME/.config)/goremote/config.yaml` is also tried when nothing is given.

Task definitions can be specified in the config file only.
Have a look to the example [config.yaml](./example/config.yaml)

### Build
Just run make
```
make
```