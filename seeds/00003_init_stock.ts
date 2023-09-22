import { Knex } from "knex";

const schema = "everytrack_backend";
const table = "stock";

export async function seed(knex: Knex): Promise<void> {
  const haveRecords = await knex(`${schema}.${table}`)
    .count("id", {
      as: "rows",
    })
    .first();

  if ((haveRecords && Number(haveRecords.rows) > 0) === false) {
    await knex(`${schema}.${table}`).del();

    const currencies: { id: string; ticker: string }[] = await knex(
      `${schema}.currency`
    ).select("id", "ticker");
    const currenciesMap: any = {};
    currencies.forEach((currency) => {
      currenciesMap[currency.ticker] = currency.id;
    });

    const countries: { id: string; code: string }[] = await knex(
      `${schema}.country`
    ).select("id", "code");
    const countriesMap: any = {};
    countries.forEach((country) => {
      countriesMap[country.code] = country.id;
    });

    await knex(`${schema}.${table}`).insert([
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "Vanguard 500 Index Fund ETF",
        ticker: "VOO",
        current_price: "410.35",
      },
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "Microsoft Corporation",
        ticker: "MSFT",
        current_price: "336.06",
      },
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "Apple Inc",
        ticker: "AAPL",
        current_price: "174.21",
      },
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "Block, Inc",
        ticker: "SQ",
        current_price: "54.15",
      },
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "Cloudflare, Inc",
        ticker: "NET",
        current_price: "64.46",
      },
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "ASML Holding N.V.",
        ticker: "ASML",
        current_price: "577.2",
      },
      {
        country_id: countriesMap["US"],
        currency_id: currenciesMap["USD"],
        name: "Mastercard Incorporated",
        ticker: "MA",
        current_price: "416.3",
      },
    ]);
  }
}
