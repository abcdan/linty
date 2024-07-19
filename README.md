# linty 🧹
A dead simple linter where you write your tests in good ol' Javascript

## Usage 🚀
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
  "verbose": true,
  "ignore": [
    "node_modules/",
    "vendor/",
    "*.test.go",
    "ignored_file.go",
    "linty.go"
  ],
  "lint": [
    {
      "type": "go",
      "regex": ".*\\.go$",
      "linter": "go.js"
    },
    {
      "type": "py",
      "regex": ".*\\.py$",
      "linter": "py.js"
    },
    {
      "type": "php",
      "regex": ".*\\.php$",
      "linter": "php.js"
    }
  ],
  "website": "https://linty.run"
}

```
This file is the configuration for linty. You can make it ignore all the `.gitignore` files and even add your own custom ignore patterns. You can also define the linting rules for different file types

Now you can add the linters for the different file types in the same directory, here's an example for Javascript:
`js.js`
```javascript
module.exports = [
  {
    name: "Unused Variable",
    description: "Checks for unused variables in JavaScript files.",
    lint: function (input) {
      const results = [];
      const lines = input.content.split("\n");
      const unusedVariableRegex = /var\s+\w+\s*=/;

      lines.forEach((line, index) => {
        if (unusedVariableRegex.test(line)) {
          results.push({
            file: input.file,
            line: index + 1,
            issue: "Unused variable",
            result: false,
          });
        }
      });

      return results;
    },
  },
  {
    name: "Console Log",
    description: "Checks for console.log statements in JavaScript files.",
    lint: function (input) {
      const results = [];
      const lines = input.content.split("\n");
      const consoleLogRegex = /console\.log\(/;

      lines.forEach((line, index) => {
        if (consoleLogRegex.test(line)) {
          results.push({
            file: input.file,
            line: index + 1,
            issue: "Console log statement",
            result: false,
          });
        }
      });

      return results;
    },
  },
];
```
4. That's it. Now whenever you push to the main branch or create a pull request, linty will run and check your code for any issues and the pipeline will fail if any issues are found.
5. Ensure your code is better 🎉

## Why? 🤔
I really wanted something I can easily add to my own project without having to learn a special linter syntax. I just want to write some Javascript that does the work. It might not be the most efficient way to do it, but it's the way I like it. It makes it easy to add custom linters for different file types and it's easy to understand.

The cool part is, you can easily fork this and make it use another language for parsing the files. Like go, python, php, etc. You can also add more linters for different file types. It's all up to you.

## Why is there Go code in the repository? 🦫
I use Go for reading through all the files and calling the Node scripts. It's just a simple way to do it. I could have used Node for everything, but I felt that Go would be faster for reading through all the files. Not sure if it is, but it gets the job done.

## Contributing 🙏
Feel free to open a PR or an issue if you have any suggestions or improvements. I'm always open to feedback. When it's a simple task, I prefer you open a PR with the changes. If it's a bigger change, it's better to open an issue first so we can discuss it.

## License 📜
It's licensed under the MIT license. You can read the full license [here](LICENSE). The TL;DR is:
- You can do whatever you want with it
- I'm not liable for anything
- If you use it, you have to include the license

## Author 🧙‍♂️
This project is created and maintained by abcdan on GitHub.
