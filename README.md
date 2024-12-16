#  Music_api #
## Автор: Трофимов Владимир ##
---
### Содержание ###
- [Изменения](#изменения)
- [Описание модуля](#описание-модуля)
---
### Изменения ###
- Добавлен интерфейс DaoEnrichment.
- Насыщение теперь добавляется в service по интерфейсу: **+возможность к расширению**.
- Добавлен интерфейс MusicService и функциональный тип CreateMusicService, который как интерфейс для всех конструкторов потомков MusicService.
- Сервис теперь добавляется в обработчики по ассоциации, а не по композиции. Добавление происходит через функцию типа CreateMusicService: **+возможность к расширению**.
- Решена проблема с экранированием спецсимволов для БД.
- Решена проблема с угрозой SQL-инъекций.


### Описание модуля ###
#### Диаграмма классов для лучшего понимания структуры модуля: ####
![1](https://github.com/Vladimir220/music_api/blob/main/pics/class_diagram.jpg)

