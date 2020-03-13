# tfe-cli

CLI client for Terraform Enterprise.

## Environment variables

* `TFE_TOKEN`: Terraform Enterprise API token
* `TFE_ORG`: Terraform Enterprise organization
* `TFE_LOG_LEVEL`: Logging level (valid values are `debug`, `info`, `warn`, `error`,
  `fatal`, `panic`)

These values can also be specified on the command line. In this case, the environment
variables are ignored.

## Management commands

By default, `tfe-cli` does not display anything if a command succeeds (unless a result
is expected, like listing the workspaces for instance). The verbosity can be adjusted
by setting the log level accordingly.

### Workspaces

#### Examples

List existing workspaces in the organization:

```bash
tfe-cli workspace list
```

Create a new workspace:

```bash
tfe-cli workspace create my-new-workspace
```

