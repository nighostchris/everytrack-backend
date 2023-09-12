import { Knex } from "knex";

const schema = "everytrack_backend";
const table = "asset_provider_account_type";

export async function seed(knex: Knex): Promise<void> {
  await knex(`${schema}.${table}`).del();

  const assetProviders: { id: string; name: string }[] = await knex(
    `${schema}.asset_provider`
  ).select("id", "name");
  const assetProvidersMap: any = {};
  assetProviders.forEach((assetProvider) => {
    assetProvidersMap[assetProvider.name] = assetProvider.id;
  });

  await knex(`${schema}.${table}`).insert([
    {
      asset_provider_id: assetProvidersMap["Chase Bank (UK)"],
      name: "Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Chase Bank (UK)"],
      name: "Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Bank of China (HK)"],
      name: "HKD Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Bank of China (HK)"],
      name: "Foreign Currency Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Bank of China (HK)"],
      name: "HKD Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Bank of China (HK)"],
      name: "USD Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (HK)"],
      name: "HKD Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (HK)"],
      name: "HKD Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (UK)"],
      name: "Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (UK)"],
      name: "Online Bonus Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (UK)"],
      name: "Regular Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (UK)"],
      name: "Flexible Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (UK)"],
      name: "Fixed Rate Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["HSBC (UK)"],
      name: "Loyalty Cash ISA",
    },
    {
      asset_provider_id: assetProvidersMap["Lloyds Bank (UK)"],
      name: "Classic Account",
    },
    {
      asset_provider_id: assetProvidersMap["Lloyds Bank (UK)"],
      name: "Easy Saver",
    },
    {
      asset_provider_id: assetProvidersMap["Lloyds Bank (UK)"],
      name: "Cash ISA Saver",
    },
    {
      asset_provider_id: assetProvidersMap["Lloyds Bank (UK)"],
      name: "1 Year Fixed Rate Cash ISA",
    },
    {
      asset_provider_id: assetProvidersMap["Citibank (HK)"],
      name: "Checking Account",
    },
    {
      asset_provider_id: assetProvidersMap["Citibank (HK)"],
      name: "HKD Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Citibank (HK)"],
      name: "USD Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Citibank (HK)"],
      name: "GBP Global Wallet",
    },
    {
      asset_provider_id: assetProvidersMap["Hang Seng Bank (HK)"],
      name: "HKD Savings / Current Wallet",
    },
    {
      asset_provider_id: assetProvidersMap["Metro Bank (UK)"],
      name: "Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Metro Bank (UK)"],
      name: "Instant Access Savings",
    },
    {
      asset_provider_id: assetProvidersMap["Metro Bank (UK)"],
      name: "Instant Access Cash ISA",
    },
    {
      asset_provider_id: assetProvidersMap["Metro Bank (UK)"],
      name: "Fixed Rate Cash ISA",
    },
    {
      asset_provider_id: assetProvidersMap["Metro Bank (UK)"],
      name: "Fixed Term Savings",
    },
    {
      asset_provider_id: assetProvidersMap["Monzo (UK)"],
      name: "Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Monzo (UK)"],
      name: "Instant Access Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Monzo (UK)"],
      name: "Easy Access Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Monzo (UK)"],
      name: "Easy Access ISA Account",
    },
    {
      asset_provider_id: assetProvidersMap["Monzo (UK)"],
      name: "Fixed Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "FlexPlus",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "FlexDirect",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "FlexAccount",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "1 Year Fixed Rate ISA",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "1 Year Triple Access Online ISA",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "Flex Instant Saver",
    },
    {
      asset_provider_id: assetProvidersMap["Nationwide (UK)"],
      name: "Instant Access Saver",
    },
    {
      asset_provider_id: assetProvidersMap["Starling (UK)"],
      name: "Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Standard Chartered (HK)"],
      name: "HKD Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Standard Chartered (HK)"],
      name: "Payroll Account",
    },
    {
      asset_provider_id: assetProvidersMap["Standard Chartered (HK)"],
      name: "HKD Marathon Savings Account",
    },
    {
      asset_provider_id: assetProvidersMap["Standard Chartered (HK)"],
      name: "USD Marathon Savings Account",
    },
  ]);
}
