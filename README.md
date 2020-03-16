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

Manage workspaces for an organization.

#### Examples

List existing workspaces in the organization:

```bash
tfe-cli workspace list
```

Create a new workspace:

```bash
tfe-cli workspace create my-new-workspace
```

### Variables

Manage variables for a workspace.

#### Examples

List exisitng variables for a specific workspace:

Create a new variable into a specific workspace:

```bash
tfe-cli variable create exisiting-workspace --var akey=a_value
```

Update an existing variable in a specific workspace:

```bash
tfe-cli variable create exisiting-workspace -f --var akey=a_nother_value
```

When creating/updating variables, several of them of can be specified at the same time:
```bash
tfe-cli variable create exisiting-workspace \
  --var akey=a_value \
  --var bkey=b_value \
  --hcl-var hclkey=hcl_value
```
