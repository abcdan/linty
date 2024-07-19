const fs = require('fs');

function parseFiles(input) {
  const content = fs.readFileSync(input, 'utf8');
  // Add your linting logic here
  // For demonstration, we'll just check if the file contains the word "FIXME"
  return !content.includes('FIXME');
}

const inputFile = process.argv[2];
const result = parseFiles(inputFile);
console.log(result);
