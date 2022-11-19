const fs = require('fs');

const contents = {
  'DATABASE_URL': process.argv[2],
  'JWT_SECRET': process.argv[3]
};

let content = 'env_variables:\n';
Object.keys(contents).forEach((key) => {
  content += `  ${key}: ${contents[key]}\n`;
});

fs.writeFile('./environ.yaml', content, (err) => {
  if (err) throw err;
  console.log('File is created successfully.');
});
