import { Knex } from "knex";

const table = "asset_provider";
const schema = "everytrack_backend";

export async function seed(knex: Knex): Promise<void> {
  await knex(`${schema}.${table}`).del();
  await knex(`${schema}.${table}`).insert([
    { name: "Chase Bank (UK)", icon: "/chase_bank_uk.svg", type: "bank" },
    { name: "HSBC (HK)", icon: "/hsbc_hk.svg", type: "bank" },
    { name: "HSBC (UK)", icon: "/hsbc_uk.svg", type: "bank" },
    { name: "Bank of China (HK)", icon: "/boc_hk.svg", type: "bank" },
    { name: "Lloyds Bank (UK)", icon: "/lloyds_uk.svg", type: "bank" },
    { name: "Citibank (HK)", icon: "/citibank_hk.svg", type: "bank" },
    { name: "Hang Seng Bank (HK)", icon: "/hang_seng_hk.svg", type: "bank" },
    { name: "Metro Bank (UK)", icon: "/metro_bank_uk.svg", type: "bank" },
    { name: "Monzo (UK)", icon: "/monzo_uk.svg", type: "bank" },
    { name: "Nationwide (UK)", icon: "/nationwide_uk.svg", type: "bank" },
    {
      name: "Standard Chartered (HK)",
      icon: "/standard_chartered_hk.svg",
      type: "bank",
    },
    { name: "Starling (UK)", icon: "/starling_uk.svg", type: "bank" },
  ]);
}
