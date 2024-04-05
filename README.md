## Deprecation: This repository is no longer supported. Please use the [Terrapak](https://github.com/eunanhardy/terrapak)

# Terrapak for Github Actions
This action integrates with your GitHub pull requests to automatically publish new versions of your Terraform modules in conjunction with [Terrapak](https://github.com/eunanhardy/terrapak)

## Getting Started

Terrapak uses a configuration file to define the modules you want to publish. Create a file named `terrapak.hcl` in the root of your repository. The file should contain a list of modules you want to publish. Each module should have a name and a path to the module directory. The path is relative to the root of the repository.
Example `terrapak.hcl` file:

```hcl
terrapak {
    hostname = "terrapak.dev"
    organization = "myorg"
}

module "aws-bucket" {
    provider = "aws"
    path = "modules/aws/bucket"
    version = "1.0.0"
    # Example url: terrapak.dev/myorg/aws-bucket/aws
}

```

Example usage as module source:
```hcl
module "bucket" {
    source = "terrapak.dev/myorg/aws-bucket/aws"
    version = "1.0.0"
    bucket_name = "my-bucket"
}
```


### Requirements
To use this action you must first setup a instance of [Terrapak](https://github.com/eunanhardy/terrapak)


### Workflows
Add the following workflow to your repository to automatically publish new versions of your modules when a pull request is merged to the target branch. 

While your pull request is open, the module is considered a draft and will accept changes willingly, but will not be permanent. Once the pull request is merged, that version of the module will be published to the registry and will no longer accept further changes.

```yaml
# terrapak_pull_request.yml
name: "Run Terrapak Sync"

on:
  pull_request:
    types: [opened, synchronize]


jobs:
  module-sync:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with: 
          fetch-depth: 0
      - name: Config-Workspace
        run: git config --global --add safe.directory /github/workspace
      - name: Terrapak Sync
        uses: eunanhardy/terrapak-action@v1
        with:
          action: sync
          github_token: ${{secrets.GITHUB_TOKEN}}
          token: ${{secrets.TP_TOKEN}}

```
Workflows for Publishing and Unpublishing modules are also available.
```yaml
# terrapak_close.yml
name: "Module Cleanup..."

on:
  pull_request:
    types: [closed]

jobs:
  module-remove:
    runs-on: ubuntu-latest
    if: github.event.pull_request.merged == false
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: "Terrapak Remove"
        uses: eunanhardy/terrapak-action@v1
        with:
          action: closed
  module-publish:
    runs-on: ubuntu-latest
    if: github.event.pull_request.merged
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: "Terrapak Publish"
        uses: eunanhardy/terrapak-action@v1
        with:
          action: merged
```
