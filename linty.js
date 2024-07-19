const fs = require("fs");
const path = require("path");

const [linterFile, fileToLint] = process.argv.slice(2);

console.log(`Linter file: ${linterFile}`);
console.log(`File to lint: ${fileToLint}`);

try {
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
} catch (error) {
  console.error(`Error running linter: ${error.message}`);
  process.exit(1);
}
