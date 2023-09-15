import { Knex } from "knex";

const schema = "everytrack_backend";

export async function seed(knex: Knex): Promise<void> {
  await knex(`${schema}.asset_provider`).insert([
    {
      name: "Futu Holdings Limited (HK)",
      icon: "/futu_hk.svg",
      type: "broker",
    },
    {
      name: "Firstrade Securities (US)",
      icon: "/firstrade_us.svg",
      type: "broker",
    },
  ]);

  const assetProviders: { id: string; name: string }[] = await knex(
    `${schema}.asset_provider`
  ).select("id", "name");
  const assetProvidersMap: any = {};
  assetProviders.forEach((assetProvider) => {
    assetProvidersMap[assetProvider.name] = assetProvider.id;
  });

  await knex(`${schema}.asset_provider_account_type`).insert([
    {
      asset_provider_id: assetProvidersMap["Futu Holdings Limited (HK)"],
      name: "Personal Brokerage Account",
    },
    {
      asset_provider_id: assetProvidersMap["Firstrade Securities (US)"],
      name: "Individual Brokerage Account",
    },
    {
      asset_provider_id: assetProvidersMap["Firstrade Securities (US)"],
      name: "Joint Brokerage Account",
    },
    {
      asset_provider_id: assetProvidersMap["Firstrade Securities (US)"],
      name: "Traditional IRA",
    },
  ]);
}
