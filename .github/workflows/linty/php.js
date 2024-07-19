module.exports = [
  {
    name: "Unused Variable",
    description: "Checks for unused variables in PHP files.",
    lint: function (input) {
      const results = [];
      const lines = input.content.split("\n");
      const unusedVariableRegex = /\$[a-zA-Z_]\w*\s*=/;

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
    name: "Echo Statement",
    description: "Checks for echo statements in PHP files.",
    lint: function (input) {
      const results = [];
      const lines = input.content.split("\n");
      const echoStatementRegex = /echo\s+/;

      lines.forEach((line, index) => {
        if (echoStatementRegex.test(line)) {
          results.push({
            file: input.file,
            line: index + 1,
            issue: "Echo statement",
            result: false,
          });
        }
      });

      return results;
    },
  },
];
