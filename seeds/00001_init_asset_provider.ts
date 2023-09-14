import { Knex } from "knex";

const table = "asset_provider";
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
      { name: "Chase Bank (UK)", icon: "/chase_bank_uk.svg", type: "savings" },
      { name: "HSBC (HK)", icon: "/hsbc_hk.svg", type: "savings" },
      { name: "HSBC (UK)", icon: "/hsbc_uk.svg", type: "savings" },
      { name: "Bank of China (HK)", icon: "/boc_hk.svg", type: "savings" },
      { name: "Lloyds Bank (UK)", icon: "/lloyds_uk.svg", type: "savings" },
      { name: "Citibank (HK)", icon: "/citibank_hk.svg", type: "savings" },
      {
        name: "Hang Seng Bank (HK)",
        icon: "/hang_seng_hk.svg",
        type: "savings",
      },
      { name: "Metro Bank (UK)", icon: "/metro_bank_uk.svg", type: "savings" },
      { name: "Monzo (UK)", icon: "/monzo_uk.svg", type: "savings" },
      { name: "Nationwide (UK)", icon: "/nationwide_uk.svg", type: "savings" },
      {
        name: "Standard Chartered (HK)",
        icon: "/standard_chartered_hk.svg",
        type: "savings",
      },
      { name: "Starling (UK)", icon: "/starling_uk.svg", type: "savings" },
    ]);
  }
}
