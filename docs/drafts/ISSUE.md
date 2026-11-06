# Generic OAuth2 client-credentials auth (e.g. Keycloak-fronted HAPI, Medplum)

## Problem

The generic provider currently supports only `bearer`, `basic`, and `none` auth types. Connecting to FHIR servers behind an OAuth2 token endpoint—such as Keycloak-fronted HAPI, Medplum, Smile CDR, Aidbox, or national health platforms—requires a provider-specific implementation. Only Cerner currently supports client credentials, and that implementation is hardwired to Cerner's specifics.

## Proposal

Add `auth.type: oauth2_client_credentials` to the generic provider, using the existing `golang.org/x/oauth2/clientcredentials` dependency:

```yaml
fhir_server:
  provider: "generic"
  base_url: "http://localhost:8080/fhir"
  auth:
    type: "oauth2_client_credentials"
    token_url: "http://localhost:8081/realms/fhir/protocol/openid-connect/token"
    client_id: "${FHIR_CLIENT_ID}"
    client_secret: "${FHIR_CLIENT_SECRET}"
    scopes: []              # optional
    endpoint_params: {}     # optional extra form params for the token request
    auth_style: "params"    # optional: params | header | auto (default auto)
```

The PR will include example configuration and the repo's first unit tests, covering token acquisition and refresh, auth styles, endpoint parameters, error handling, configuration validation, and regression coverage for existing auth types.

Would maintainers prefer this support on the generic provider, or a dedicated provider?
