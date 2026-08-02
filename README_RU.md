# NTS Radio RSS / JSON Feed Generator

## Содержание

- [NTS Radio RSS / JSON Feed Generator](#nts-radio-rss--json-feed-generator)
  - [Содержание](#содержание)
  - [Описание](#описание)
  - [Конфигурация](#конфигурация)
  - [Структура фидов](#структура-фидов)
    - [Разделы](#разделы)
    - [Структура ответа](#структура-ответа)
      - [1. RSS](#1-rss)
      - [2. JSON](#2-json)
    - [URL Query параметры запросов](#url-query-параметры-запросов)
  - [Структура запросов API](#структура-запросов-api)
    - [1. Информация об эпизоде](#1-информация-об-эпизоде)
    - [2. Треклист эпизода](#2-треклист-эпизода)
  - [Типы данных](#типы-данных)
    - [Элемент треклиста (Tracklist Item)](#элемент-треклиста-tracklist-item)
    - [Эпизод (Episode)](#эпизод-episode)
    - [Жанр (Genre)](#жанр-genre)
    - [Медиа (Media)](#медиа-media)
    - [Источник (Source)](#источник-source)
  - [Коды ошибок](#коды-ошибок)

## Описание

Веб-сервер для генерации RSS / JSON фидов NTS Radio с возможностью управления настройками через консоль командной строки Linux. Поддерживает загрузку треклистов по ссылке.

## Конфигурация

```yaml
nts:
  api_url: "https://www.nts.live/api/v2/"

api:
  public_url: "https://nts-feed/api" # Публичный адрес сервера

tracklist:
  cache_ttl: "1h"
  enabled: true

feed:
  output_dir: "./feeds"
  include_tracklist_url: true
```

## Структура фидов

### Разделы

| Раздел  | Метод | Путь                            |
| ------- | ----- | ------------------------------- |
| Latest  | GET   | `/feeds/latest[?query]`         |
| Подкаст | GET   | `/feeds/show/{show_id}[?query]` |

### Структура ответа

#### 1. RSS

```xml
<channel>
  <title>Channel Title</title>
  <link>https://www.nts.live/{channel_path}</link>
  <description>Channel description</description>
  <category>Channel category</category>
  <image>Channel Image Url</image>
  <lastBuildDate>Last Updated Date</lastBuildDate>
  <item>
      <title>Episode Title</title>
      <link>https://www.nts.live/{episode_path}</link>
      <description>Episode description</description>
      <pubDate>2024-01-15T20:00:00Z</pubDate>
      <enclosure url="{audio_url}" type="audio/mpeg" length="0"/>
      <guid>episode-123</guid>

      <nts:show_id>Show Id</nts:show_id>
      <nts:air_date>2024-01-15</nts:air_date>
      <nts:duration>120</nts:duration>

      <api:tracklist_url>Api Tracklist Url</api:tracklist_url>
  </item>
</channel>
```

#### 2. JSON

```json
{
  "channel": {
    "title": "Channel Title",
    "link": "https://www.nts.live/{channel_path}",
    "description": "Channel description",
    "category": "Channel category",
    "image": "Channel Image Url",
    "lastBuildDate": "Last Updated Date",
    "item": {
      "title": "Episode Title",
      "link": "https://www.nts.live/{episode_path}",
      "description": "Episode description",
      "pubDate": "2024-01-15T20:00:00Z",
      "enclosure": {
        "url": "{audio_url}",
        "type": "audio/mpeg",
        "length": "0"
      },
      "guid": "episode-123",
      "nts:show_id": "Show Id",
      "nts:air_date": "2024-01-15",
      "nts:duration": "120",
      "tracklist": [
        {
          "artist": "string",
          "title": "string",
          "uid": "string (UUID or null)",
          "offset": "integer or null",
          "duration": "integer or null",
          "offset_estimate": "integer or null",
          "duration_estimate": "integer or null",
          "acr_id": "string or null",
          "deezer_track_id": "integer or null",
          "isrc_id": "string or null",
          "musicbrainz_track_id": "string (UUID or null)"
        }
      ]
    }
  }
}
```

### URL Query параметры запросов

| Параметр   | Тип      | По умолчанию | Описание                                   | Пример                         |
| ---------- | -------- | ------------ | ------------------------------------------ | ------------------------------ |
| `format`   | string   | `rss`        | Формат ответа: `rss`, `atom`, `json`       | `?format=json`                 |
| `limit`    | integer  | `50`         | Максимум элементов на страницу (1–200)     | `?limit=20`                    |
| `offset`   | integer  | `0`          | Пропуск элементов для пагинации            | `?offset=40`                   |
| `page`     | integer  | `1`          | Номер страницы (альтернатива `offset`)     | `?page=3`                      |
| `since`    | datetime | —            | Элементы позже указанной даты (ISO 8601)   | `?since=2026-08-01T00:00:00Z`  |
| `before`   | datetime | —            | Элементы раньше указанной даты (ISO 8601)  | `?before=2026-08-01T00:00:00Z` |
| `category` | string   | —            | Фильтр по категории / тегу                 | `?category=tech`               |
| `author`   | string   | —            | Фильтр по имени пользователя / ID автора   | `?author=johndoe`              |
| `q`        | string   | —            | Поиск по заголовкам и содержимому          | `?q=golang+tutorial`           |
| `order`    | string   | `desc`       | Порядок сортировки: `desc` или `asc`       | `?order=asc`                   |
| `full`     | boolean  | `false`      | Полное содержимое вместо краткого описания | `?full=true`                   |
| `callback` | string   | —            | Имя функции JSONP-callback                 | `?callback=myCallback`         |

## Структура запросов API

### 1. Информация об эпизоде

`GET /shows/{show_id}/episodes/{episode_id}`

Возвращает [Эпизод](#эпизод-episode) и вложенный список [Элемент треклиста](#элемент-треклиста-tracklist-item):

```json
{
  "embeds": {
    "tracklist": {
      "metadata": {
        "resultset": {
          "count": "integer",
          "offset": "integer",
          "limit": "integer"
        }
      },
      "results": ["[Элемент треклиста]"]
    }
  }
}
```

### 2. Треклист эпизода

`GET /shows/{show_id}/episodes/{episode_id}/tracklist`

Возвращает список [Элемент треклиста](#элемент-треклиста-tracklist-item)

```json
{
  "metadata": {
    "resultset": { "count": "integer", "offset": "integer", "limit": "integer" }
  },
  "results": ["[Элемент треклиста](#элемент-треклиста-tracklist-item)"]
}
```

## Типы данных

### Элемент треклиста (Tracklist Item)

```json
{
  "artist": "string",
  "title": "string",
  "uid": "string (UUID) or null",
  "offset": "integer or null",
  "duration": "integer or null",
  "offset_estimate": "integer or null",
  "duration_estimate": "integer or null",
  "acr_id": "string or null",
  "deezer_track_id": "integer or null",
  "isrc_id": "string or null",
  "musicbrainz_track_id": "string (UUID) or null"
}
```

### Эпизод (Episode)

```json
{
  "status": "string",
  "updated": "datetime (ISO 8601)",
  "name": "string",
  "description": "string",
  "description_html": "string (HTML)",
  "external_links": "array",
  "moods": "array",
  "genres": "Genre[]",
  "location_short": "string",
  "location_long": "string",
  "intensity": "string (numeric)",
  "media": "Media",
  "episode_alias": "string",
  "show_alias": "string",
  "broadcast": "datetime (ISO 8601)",
  "mixcloud": "string (URL)",
  "audio_sources": "AudioSource[]",
  "brand": "object"
}
```

### Жанр (Genre)

```json
{
  "id": "string",
  "value": "string"
}
```

### Медиа (Media)

```json
{
  "background_large": "string (URL)",
  "background_medium_large": "string (URL)",
  "background_medium": "string (URL)",
  "background_small": "string (URL)",
  "background_thumb": "string (URL)",
  "picture_large": "string (URL)",
  "picture_medium_large": "string (URL)",
  "picture_medium": "string (URL)",
  "picture_small": "string (URL)",
  "picture_thumb": "string (URL)"
}
```

### Источник (Source)

```json
{
  "url": "string (URL)",
  "source": "string"
}
```

## Коды ошибок

| Код | Статус                | Описание                           |
| --- | --------------------- | ---------------------------------- |
| 200 | OK                    | Успешный запрос                    |
| 304 | Not Modified          | Фид не изменился (совпадение ETag) |
| 400 | Bad Request           | Некорректные параметры запроса     |
| 404 | Not Found             | Фид или элементы не найдены        |
| 429 | Too Many Requests     | Превышен лимит запросов            |
| 500 | Internal Server Error | Внутренняя ошибка сервера          |
