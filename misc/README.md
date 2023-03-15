## Настройка окружения для разработки

### Разработка в VS Code под sudo на удаленном сервере

1. Скопируйте `settings.json` и `launch.json` в папку `.vscode` в корне проекта.
2. В `launch.json` замените переменные окружения `WG_GRPC_API_ADDRESS`, `WG_GRPC_API_ENDPOINT` и `WG_GRPC_API_PEER_FOLDER`.
3. Теперь `wg-grpc-api` можно отлаживать на удаленном сервере под `sudo`.


Если для запуска проекта в режиме отладки требуется пароль sudo,
то скрипт `dlv.sh` необходимо поменять на:
```sh
#!/bin/bash
echo "dlv wrapper script to run debugger as sudo"
if [ "$DEBUG_AS_ROOT" = "true" ]; then
	exec echo <SUDO_PASSWORD> | /usr/bin/sudo -S -E $HOME/go/bin/dlv --only-same-user=false "$@"
else
	exec $HOME/go/bin/dlv "$@"
fi
```
