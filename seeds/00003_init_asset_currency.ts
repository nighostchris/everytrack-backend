import { Knex } from "knex";

const table = "currency";
const schema = "everytrack_backend";

export async function seed(knex: Knex): Promise<void> {
  await knex(`${schema}.${table}`).del();
  await knex(`${schema}.${table}`).insert([
    { name: "HKD", symbol: "HKD$" },
    { name: "USD", symbol: "USD$" },
    { name: "GBP", symbol: "£" },
  ]);
}
