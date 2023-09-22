import { Knex } from "knex";

const table = "country";
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
      { name: "Hong Kong", code: "HK" },
      { name: "United States", code: "US" },
      { name: "United Kingdom", code: "UK" },
    ]);
  }
}
