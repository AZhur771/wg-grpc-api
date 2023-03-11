[![Go](https://github.com/AZhur771/wg-grpc-api/actions/workflows/ci.yaml/badge.svg)](https://github.com/AZhur771/wg-grpc-api/actions/workflows/ci.yaml)

# wg-grpc-api

GRPC API to manage wireguard peers.

## Environment variables
| Variable                    | Default       | Comment                                                          |
|:----------------------------|:--------------|:-----------------------------------------------------------------|
| WG_GRPC_API_HOST            | localhost     | App host                                                         |
| WG_GRPC_API_PORT            | 3000          | Grpc api port                                                    |
| WG_GRPC_API_GATEWAY         | false         | Enale grpc api gateway                                           |
| WG_GRPC_API_GATEWAY_PORT    | 3001          | Grpc api gateway port                                            |
| WG_GRPC_API_SWAGGER         | false         | Enable grpc api gateway swagger docs (served at /swagger-ui)     |
| WG_GRPC_API_DEVICE          | wg0           | Wireguard interface                                              |
| WG_GRPC_API_ADDRESS         | -             | Wireguard virtual address in CIDR notation (required)            |
| WG_GRPC_API_ENDPOINT        | -             | VPN server public ip (required)                                  |
| WG_GRPC_API_TOKENS          | -             | Authentication tokens                                            |
| WG_GRPC_API_SERVER_CERT     | -             | Server certificate                                               |
| WG_GRPC_API_SERVER_KEY      | 6379          | Server key                                                       |
| WG_GRPC_API_PEER_FOLDER     | -             | Path to folder containing configured peers                       |


## TODO
- [ ] Unit tests
- [ ] Converting peer conf files to json
- [ ] Fix enable/disable functionality
- [ ] Add docker compose
