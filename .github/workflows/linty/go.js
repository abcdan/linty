module.exports = [
  {
    name: "Unused Import",
    description: "Checks for unused imports in Go files.",
    lint: function (input) {
      const results = [];
      const lines = input.split("\n");
      const unusedImportRegex = /import\s+\w+\s+"[^"]+"/;

      lines.forEach((line, index) => {
        if (unusedImportRegex.test(line)) {
          results.push({
            file: input.file,
            line: index + 1,
            issue: "Unused import",
            result: false,
          });
        }
      });

      return results;
    },
  },
  {
    name: "Unused Variable",
    description: "Checks for unused variables in Go files.",
    lint: function (input) {
      const results = [];
      const lines = input.split("\n");
      const unusedVariableRegex = /var\s+\w+\s+/;

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
    name: "Missing Error Handling",
    description: "Checks for missing error handling in Go files.",
    lint: function (input) {
      const results = [];
      const lines = input.split("\n");
      const errorHandlingRegex = /if\s+err\s*!=\s*nil\s*{\s*return\s+err\s*}/;

      lines.forEach((line, index) => {
        if (!errorHandlingRegex.test(line) && line.includes("if err != nil")) {
          results.push({
            file: input.file,
            line: index + 1,
            issue: "Missing error handling",
            result: false,
          });
        }
      });

      return results;
    },
  },
];
