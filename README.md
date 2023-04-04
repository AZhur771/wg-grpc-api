[![Go](https://github.com/AZhur771/wg-grpc-api/actions/workflows/ci.yaml/badge.svg)](https://github.com/AZhur771/wg-grpc-api/actions/workflows/ci.yaml)

# wg-grpc-api

GRPC and Rest API to manage wireguard peers.

## Supported environment variables
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
| WG_GRPC_API_SERVER_KEY      | -             | Server key                                                       |
| WG_GRPC_API_PEER_FOLDER     | -             | Path to folder containing configured peers                       |


## TODO
- [x] Fix enable/disable functionality.
- [x] Add unit tests.
- [ ] Allow modifying more server options (PreUp/PostUp/PreDown/PostDown/DNS/MTU/ListenPort/Table).
