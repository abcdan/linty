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
