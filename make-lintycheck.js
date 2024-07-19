const fs = require("fs");
const crypto = require("crypto");

function calculateMD5(file) {
  const fileContent = fs.readFileSync(file);
  const md5 = crypto.createHash("md5");
  md5.update(fileContent);
  return md5.digest("hex");
}

function generateLintycheckFile() {
  const lintyGoChecksum = calculateMD5("linty.go");
  const lintyJsChecksum = calculateMD5("linty.js");

  const lintycheckContent = `linty.go|${lintyGoChecksum}
linty.js|${lintyJsChecksum}`;

  fs.writeFileSync("LINTYCHECK", lintycheckContent);
  console.log("LINTYCHECK file generated successfully.");
}

generateLintycheckFile();
