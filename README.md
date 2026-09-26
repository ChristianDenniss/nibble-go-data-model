# go-data-model

The domain of the system: what a restaurant, menu item, offer, or price observation *is*, what persistence must be able to do, and the business rules for changing that state.

**Why this repo exists:** those rules should not live in the API process or in SQL. If `data-acquisition` or a second API needs to create a restaurant, it should call the same service, not copy structs. This module is the inner hexagon. It does not import Postgres, gRPC, or HTTP.

Each context is a folder that can grow without turning into a dump:

```text
restaurant/
  entity/        types and domain errors
  repository/    persistence contract the domain requires
  service/       business logic
```

`money` and `location` are value objects. They have no repository.

Services should pin **`v1.1.0` or later**. Tag `v1.0.0` is the pre-split scaffold
and does not contain these packages.

Storage adapters (see `go-data-store`) implement the repository interfaces. Transports (see `api-engine`) call the services.
