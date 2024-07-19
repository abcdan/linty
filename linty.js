const fs = require("fs");
const path = require("path");

const [linterFile, fileToLint] = process.argv.slice(2);

const linterConfig = require(path.join(__dirname, linterFile));

const fileContent = fs.readFileSync(fileToLint, "utf-8");

function runLint(input) {
  const results = [];
  linterConfig.forEach((linter) => {
    results.push(...linter.lint(input));
  });
  return results;
}

const lintResults = runLint({ file: fileToLint, content: fileContent });
console.log(JSON.stringify(lintResults));
