const fs = require("fs");
const path = require("path");

const [linterFile, fileToLint] = process.argv.slice(2);

const configDir = path.dirname(process.argv[1]);

const linterConfig = require(path.join(configDir, linterFile));

const fileContent = fs.readFileSync(fileToLint, "utf-8");

function runLint(input) {
  const results = [];
  linterConfig.forEach((linter) => {
    try {
      results.push(...linter.lint(input));
    } catch (error) {
      console.error(`Error running linter '${linter.name}': ${error.message}`);
      results.push({
        file: input.file,
        result: false,
        error: `Linter '${linter.name}' failed: ${error.message}`,
      });
    }
  });
  return results;
}

const lintResults = runLint({ file: fileToLint, content: fileContent });
console.log(JSON.stringify(lintResults));
