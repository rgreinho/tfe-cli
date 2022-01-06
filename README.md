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
* `TFE_ADDRESS`: Terraform Enterprise API address
* `TFE_BASEPATH`: Base path on which the Terraform Enterprise API is served.

Some of these values can also be specified on the command line. In this case, the
environment variables are ignored.

## Management commands

By default, `tfe-cli` does not display anything if a command succeeds (unless a result
is expected, like listing the workspaces for instance). The verbosity can be adjusted
by setting the log level accordingly.

### Workspaces

Manage workspaces for an organization.

#### Create

Create a new TFE workspace.

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

##### Example

```bash
tfe-cli workspace delete my-new-workspace
```

#### List

List existing workspaces in the organization.

##### Example

```bash
tfe-cli workspace list
```

### Variables

Manage variables for a workspace.

#### Create

##### Examples

Create a new variable into a specific workspace:

```bash
tfe-cli variable create exisiting-workspace --var akey=a_value
```

Update an existing variable in a specific workspace:

```bash
tfe-cli variable create my-exisiting-workspace -f --var akey=another_value
```

When creating/updating variables, several of them of can be specified at the
same time:

```bash
tfe-cli variable create my-exisiting-workspace \
  --var akey=a_value \
  --svar bkey=b_secure_value \
  --hvar hclkey=hcl_value \
  --var-file stage.tfvars \
```

#### Delete

##### Example

Delete a variable:

```bash
tfe-cli variable delete my-workspace backend_port
```

#### List

List exisitng variables for a specific workspace.

##### Example

List variables:

```bash
tfe-cli variable list my-workspace
```

### Notifications

#### List

List TFE notifications for a specific workspace.

##### Example

```bash
tfe-cli notification list my-workspace
```

#### Create

Creates or update a notification.

##### Example

Create a Slack notification for the `created` and `errored` events:

```bash
tfe-cli notification create my-workspace \
  my-notification \
  --type slack \
  --url https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX \
  --triggers run:created \
  --triggers run:errored
```

#### Delete

Deletes a notification by its name, in a specific workspace.

##### Example

```bash
tfe-cli notification delete my-workspace my-notification
```
