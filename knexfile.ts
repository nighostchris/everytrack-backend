import "dotenv/config";
import { Knex } from "knex";

const DATABASE = process.env.DATABASE;

const config: Knex.Config = {
  client: "pg",
  connection: DATABASE,
  seeds: {
    extension: "ts",
    directory: "./seeds",
  },
};

export default config;
