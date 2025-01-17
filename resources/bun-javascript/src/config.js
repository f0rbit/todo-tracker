// load config from file
import args from "./args";

// read using bun file
const file = Bun.file(`${args.config}`);
const config = await file.json();

if (!config.tags) throw new Error("No tags found in config file");

const { tags, ignore = [] } = config;

// export
export default { tags, ignore };
