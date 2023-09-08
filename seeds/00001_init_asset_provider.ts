import { Knex } from "knex";

const table = "asset_provider";
const schema = "everytrack_backend";

export async function seed(knex: Knex): Promise<void> {
  await knex(`${schema}.${table}`).del();
  await knex(`${schema}.${table}`).insert([
    { name: "Chase Bank UK", icon: "/chase_bank_uk.svg" },
  ]);
}
