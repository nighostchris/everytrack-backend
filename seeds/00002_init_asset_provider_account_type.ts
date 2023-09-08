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
      asset_provider_id: assetProvidersMap["Chase Bank UK"],
      name: "Current Account",
    },
    {
      asset_provider_id: assetProvidersMap["Chase Bank UK"],
      name: "Savings Account",
    },
  ]);
}
