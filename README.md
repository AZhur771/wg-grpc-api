# wg-grpc-api

GRPC API to manage wireguard peers.

## Supported environment variables
| Variable                    | Default       | Comment                                                          |
|:----------------------------|:--------------|:-----------------------------------------------------------------|
| WG_GRPC_API_PRODUCTION      | false         | App mode                                                         |
| WG_GRPC_API_HOST            | localhost     | App host                                                         |
| WG_GRPC_API_PORT            | 3000          | Grpc api port                                                    |
| WG_GRPC_API_GATEWAY         | false         | Enale grpc api gateway                                           |
| WG_GRPC_API_GATEWAY_PORT    | 3001          | Grpc api gateway port                                            |
| WG_GRPC_API_SWAGGER         | false         | Enable grpc api gateway swagger docs (served at /swagger-ui)     |
| WG_GRPC_API_DEVICE          | wg0           | Wireguard interface                                              |
| WG_GRPC_API_ADDRESS         | -             | Wireguard virtual address in CIDR notation (required)            |
| WG_GRPC_API_ENDPOINT        | -             | VPN server public ip (required)                                  |
| WG_GRPC_API_REDIS_HOST      | localhost     | Redis host                                                       |
| WG_GRPC_API_REDIS_PORT      | 6379          | Redis port                                                       |
| WG_GRPC_API_REDIS_PASSWORD  | -             | Redis password                                                   |


## TODO
- [ ] Integration tests
- [ ] Syncing existing peers with redis
- [ ] Fix enable/disable functionality
- [ ] CLI (grpc client)
