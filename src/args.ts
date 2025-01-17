// parse args from command line
import { parseArgs } from "util";

const { values, positionals } = parseArgs({
  args: Bun.argv,
  options: {
    dir: {
      type: 'string',
    },
    config: {
      type: 'string',
    },
  },
  strict: true,
  allowPositionals: true,
});

const args = { dir: values.dir, config: values.config, path: positionals[1] };

export default args;
