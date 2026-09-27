# nibble-go-data-model

The domain of Nibble: what a restaurant, menu item, offer, or price observation *is*, what persistence must be able to do, and the business rules for changing that state.

**Why this repo exists:** those rules should not live in the API process or in SQL. If `nibble-data-acquisition` or a second API needs to create a restaurant, it should call the same service, not copy structs. This module is the inner hexagon. It does not import Postgres, gRPC, or HTTP.

Each context is a folder that can grow without turning into a dump:

```text
restaurant/
  entity/        types and domain errors
  repository/    persistence contract the domain requires
  service/       business logic
```

Contexts (legacy, still wired for ingest v1): `restaurant`, `menu`, `offer`, `observation`, `provider`, `category`, `cuisine`, `account`, `cart`, `order`.

Target contexts: `channel`, `ingest`, `source`, `brand`, `place`, `dish`, `resolution`, `itemprice`, `quoteobs`, `promotion`, `market`, `serviceability`, `user`, `compare`, `merchandising` (sponsored placements), `home` (home feed read model).

Read models: `storefront` (legacy catalog bundle for the web browse experience).

`money`, `location`, and `restaurant/entity.Rating` are value objects. They have no repository.

TypeScript for the browser is generated from these structs (`go run ./cmd/gentypes`). Do not redefine them in `nibble-web-platform`.

Services should pin **`v1.1.0` or later**. Tag `v1.0.0` is the pre-split scaffold
and does not contain these packages.

Storage adapters (see `nibble-go-data-store`) implement the repository interfaces. Transports (see `nibble-api-engine`) call the services.

## Provider catalog

`catalog` validates acquisition bundles and builds the public restaurant/menu
projection. `catalog/branch_matches.json` contains reviewed branch evidence;
item normalization removes only cosmetic trademark/calorie/whitespace changes,
retaining size, quantity and meal wording. Missing provider prices remain null.
The public types deliberately omit collection provenance.
