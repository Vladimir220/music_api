#  Music_api #
## Автор: Трофимов Владимир ##

## Внимение! ##
Для работы дополнительных обогатителей daoLastFm и daoLyricsCom нужны переменные окружения ("TOKEN_LASTFM", "TOKEN_LYRICSCOM", "UID_LYRICSCOM"), которые я из соображений безопасности не выложил, поэтому по цепочке обогатителей будет срабатывать только trackEnricherDefault.

## Docker ##
Для запуска системы через docker-compose необходимо внести следующие изменения...

В .env: 
- укажите адрес сервера обогатителя, если есть (ENRCH_SERVER_HOST)
- установите DOCKER_MOD в 1

    ```
    ENRCH_SERVER_HOST=""
    DOCKER_MOD="1"
    ```
    