// load config from file
import args from "./args";

// read using bun file
const file = Bun.file(`${args.dir}/${args.config}`);
const config = await file.json();

// export
export default config;
