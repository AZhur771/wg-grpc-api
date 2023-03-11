## Настройка окружения для разработки

### Разработка в VS Code под sudo на удаленном сервере

1. Скопируйте `settings.json` и `launch.json` в папку `.vscode` в корне проекта.
2. В `launch.json` замените переменные окружения `WG_GRPC_API_ADDRESS`, `WG_GRPC_API_ENDPOINT` и `WG_GRPC_API_PEER_FOLDER`.
3. Теперь `wg-grpc-api` можно отлаживать на удаленном сервере под `sudo`.
