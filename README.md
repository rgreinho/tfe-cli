# tfe-cli

CLI client for Terraform Enterprise.

## Installation

For Linux/OSX, run the following command from a terminal to get the latest version:

```bash
bash <(curl -sSfL https://raw.githubusercontent.com/rgreinho/tfe-cli/master/extras/tfe-cli-installer.sh)
```

For Windows, download the binary from the [release page](https://github.com/rgreinho/tfe-cli/releases).

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

```bash
$ tfe-cli workspace
Manage TFE workspaces.

Usage:
  tfe-cli workspace [command]

Available Commands:
  create      Create a TFE workspace
  delete      Delete a TFE workspace
  list        List TFE workspaces
```

#### Create

Create a new TFE workspace.

```bash
Usage:
  tfe-cli workspace create [WORKSPACE] [flags]

Flags:
      --autoapply                 Apply changes automatically
      --filetriggers              Filter runs based on the changed files in a VCS push
  -f, --force                     Update workspace if it exists
  -h, --help                      help for create
      --terraformversion string   Specify the Terraform version
      --vcsrepository string      Specify a workspace's VCS repository
      --workingdirectory string   Specify a relative path that Terraform will execute within
```

The format of the VCS option is string of colon sperated values: `<OAuthTokenID>:<repository>:<branch>`.

##### Examples

Create a new workspace with default values:

```bash
tfe-cli workspace create my-new-workspace
```

Setup the VCS Repository:

```bash
tfe-cli workspace create my-new-workspace --vcsrepository ot-8Xc1NTYpjIQZIwIh:organization/repository:master
```

#### Delete

Delete an exisiting workspace.

```bash
Usage:
  tfe-cli workspace delete [WORKSPACE]
```

#### Example

```bash
tfe-cli workspace delete my-new-workspace
```

#### List

List existing workspaces in the organization.

```bash
Usage:
  tfe-cli workspace list
```

### Variables

Manage variables for a workspace.

```bash
tfe-cli variable
Manage TFE variables.

Usage:
  tfe-cli variable [command]

Available Commands:
  create      Create TFE variables
  delete      Delete a TFE variable for a specific workspace
  list        List TFE variables for a specific workspace
```

#### Create

```bash
Usage:
  tfe-cli variable create [WORKSPACE] [flags]

Flags:
      --evar stringArray    Create an environment variable
  -f, --force               Overwrite a variable if it exists
  -h, --help                help for create
      --hvar stringArray    Create an HCL variable
      --sevar stringArray   Create a sensitive environment variable
      --shvar stringArray   Create a sensitive HCL variable
      --svar stringArray    Create a regular sensitive variable
      --var stringArray     Create a regular variable
      --var-file string     Create non-sensitive regular and HCL variables from a file
```

##### Examples

Create a new variable into a specific workspace:

```bash
tfe-cli variable create exisiting-workspace --var akey=a_value
```

Update an existing variable in a specific workspace:

```bash
tfe-cli variable create my-exisiting-workspace -f --var akey=another_value
```

When creating/updating variables, several of them of can be specified at the same time:
```bash
tfe-cli variable create my-exisiting-workspace \
  --var akey=a_value \
  --svar bkey=b_secure_value \
  --hvar hclkey=hcl_value \
  --var-file stage.tfvars \
```

#### Delete

```bash
Usage:
  tfe-cli variable delete [WORKSPACE] [VARIABLE]
```

#### Example

Delete variable:

```bash
tfe-cli variable delete my-workspace backend_port
```

#### List

List exisitng variables for a specific workspace.

```bash
Usage:
  tfe-cli variable list [WORKSPACE]
```

#### Example

List variables:

```bash
tfe-cli variable list my-workspace
```
