import { Knex } from "knex";

const table = "currency";
const schema = "everytrack_backend";

export async function seed(knex: Knex): Promise<void> {
  const haveRecords = await knex(`${schema}.${table}`)
    .count("id", {
      as: "rows",
    })
    .first();

  if ((haveRecords && Number(haveRecords.rows) > 0) === false) {
    await knex(`${schema}.${table}`).del();
    await knex(`${schema}.${table}`).insert([
      { ticker: "HKD", symbol: "HKD$" },
      { ticker: "USD", symbol: "USD$" },
      { ticker: "GBP", symbol: "£" },
    ]);
  }
}
