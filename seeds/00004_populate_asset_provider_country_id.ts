import { Knex } from "knex";

const schema = "everytrack_backend";
const table = "asset_provider";

export async function seed(knex: Knex): Promise<void> {
  const haveRecords = await knex(`${schema}.${table}`)
    .count("id", {
      as: "rows",
    })
    .first();

  if ((haveRecords && Number(haveRecords.rows) > 0) === false) {
    await knex(`${schema}.${table}`).del();

    const countries: { id: string; code: string }[] = await knex(
      `${schema}.country`
    ).select("id", "code");
    const countriesMap: any = {};
    countries.forEach((country) => {
      countriesMap[country.code] = country.id;
    });

    await knex(`${schema}.${table}`).insert([
      {
        name: "Chase Bank",
        icon: "/chase_bank_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "HSBC",
        icon: "/hsbc_hk.svg",
        type: "savings",
        country_id: countriesMap["HK"],
      },
      {
        name: "HSBC",
        icon: "/hsbc_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "Bank of China",
        icon: "/boc_hk.svg",
        type: "savings",
        country_id: countriesMap["HK"],
      },
      {
        name: "Lloyds Bank",
        icon: "/lloyds_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "Citibank",
        icon: "/citibank_hk.svg",
        type: "savings",
        country_id: countriesMap["HK"],
      },
      {
        name: "Hang Seng Bank",
        icon: "/hang_seng_hk.svg",
        type: "savings",
        country_id: countriesMap["HK"],
      },
      {
        name: "Metro Bank",
        icon: "/metro_bank_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "Monzo",
        icon: "/monzo_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "Nationwide",
        icon: "/nationwide_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "Standard Chartered",
        icon: "/standard_chartered_hk.svg",
        type: "savings",
        country_id: countriesMap["HK"],
      },
      {
        name: "Starling",
        icon: "/starling_uk.svg",
        type: "savings",
        country_id: countriesMap["UK"],
      },
      {
        name: "Futu Holdings Limited",
        icon: "/futu_hk.svg",
        type: "broker",
        country_id: countriesMap["HK"],
      },
      {
        name: "Firstrade Securities",
        icon: "/firstrade_us.svg",
        type: "broker",
        country_id: countriesMap["US"],
      },
    ]);
  }
}
