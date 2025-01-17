// load config from file
import args from "./args";
import { CONFIG_SCHEMA } from "./schema";

// read using bun file
const file = Bun.file(`${args.config}`);
const config = await file.json();

// validate with zod CONFIG_SCHEMA
const result = CONFIG_SCHEMA.safeParse(config);
if (!result.success) throw new Error(result.error.message);

const { tags, ignore = [] } = result.data;

// export
export default { tags, ignore };
