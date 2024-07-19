# linty
A dead simple linter where you write your tests in good ol' Javascript

## Usage
1. Simply create a `.github/workflows/linty.yml` file in your repository with the following content:
```yaml
name: Linty

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    - name: Run Linty
      uses: abcdan/linty@v0.0.7
      with:
        go-version: '1.21.4'
        node-version: '18'
        config-path: '.github/workflows/linty'
```
2. Create a `.github/workflows/linty` directory where you can define the following files:
`linty.json`
```json
{
  "abcdan": "linty",
  "gitignore": true,
  "verbose": false,
  "ignore": [
    "node_modules/",
    "vendor/",
    "*.test.go",
    "ignored_file.go",
    "linty.go"
  ],
  "lint": [
    { "type": "go", "regex": ".*\\.go$" },
    { "type": "py", "regex": ".*\\.py$" },
    { "type": "php", "regex": ".*\\.php$" }
  ],
  "website": "https://linty.run"
}
```
This file is the configuration for linty. You can make it ignore all the `.gitignore` files and even add your own custom ignore patterns. You can also define the linting rules for different file types.
