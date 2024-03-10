// load config from file
const DIR = `../resources`;
const FILENAME = `config.json`;

// read using bun file
const file = Bun.file(`${DIR}/${FILENAME}`);
const config = await file.json();

// export
export default config;
