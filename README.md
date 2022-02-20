# Sailing Assist Mie API Server
視覚障がい者帆走支援アプリのAPIサーバー

© 2021 NPO法人セイラビリティ三重

# API
## デバイス一覧を確認する
### Endpoint
``GET`` /devices

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## デバイス情報を確認する
### Endpoint
``GET`` /device/`:device_id`

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## デバイス情報を登録する
### Endpoint
``POST`` /devices

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## デバイス情報を更新する
### Endpoint
``PUT`` /device/`:device_id`

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## デバイス情報を削除する
### Endpoint
``DELETE`` /device/`:device_id`

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## ユーザー一覧を確認する
### Endpoint
``GET`` /users

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## ユーザー情報を確認する
### Endpoint
``GET`` /user/`:user_id`

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## ユーザー情報を登録する
### Endpoint
``POST`` /users

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## ユーザー情報を更新する
### Endpoint
``PUT`` /user/':user_id`

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## ユーザー情報を削除する
### Endpoint
``DELETE`` /user/`:user_id`

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

## レース情報を登録する
### Endpoint
``POST`` /races

### Headers
| Key | Value |
| --- | ----- |
| Authorization | Bearer `ACCESS_TOKEN` |

# データベース
## devices (デバイス情報)
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | char(16) | NOT NULL PRIMARY KEY | Android ID |
| name | varchar(32) | NOT NULL | 端末名 |
| model | smallint | NOT NULL | モデル番号 |
| latitude | double precision | | 緯度 |
| longitude | double precision | | 経度 |

CREATE TABLE devices (id char(16) NOT NULL PRIMARY KEY, name varchar(32) NOT NULL, model smallint NOT NULL, latitude double precision, longitude double precision);

## users (ユーザー情報)
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | uuid | PRIMARY KEY DEFAULT uuid_generate_v4() | ユーザーID |
| login_id | varchar(16) | NOT NULL | ログインID |
| display_name | varchar(32) | NOT NULL | 名前 |
| password | bytea | NOT NULL | パスワード |
| group_id | uuid | NOT NULL | グループID |
| role | varchar(16) | NOT NULL | `athlete`, `admin` or 'developer'
| device_id | char(16) | | デバイスID |
| sail_num | smallint | | セイル番号 |
| course_limit | float | | コースリミット |
| image_url | varchar(512) | | プロフィール画像のURL (Cloudinary) |
| note | text | | 備考 |

CREATE TABLE users (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), login_id varchar(16) NOT NULL, display_name varchar(32) NOT NULL, password bytea NOT NULL, group_id uuid NOT NULL, role varchar(16) NOT NULL, device_id char(16), sail_num smallint, course_limit float, image_url varchar(512), note text);

## races (レース情報)
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | uuid | PRIMARY KEY DEFAULT uuid_generate_v4() | レースID |
| name | varchar(32) | NOT NULL | レース名 |
| start_at | timestamp | NOT NULL | 開始時刻 |
| end_at | timestamp | NOT NULL | 終了時刻 |
| point_a | char(16) | | A地点のデバイスID |
| point_b | char(16) | | B地点のデバイスID |
| point_c | char(16) | | C地点のデバイスID |
| athlete | uuid[] | | 競技者(ユーザー)ID一覧 |
| memo | text | | メモ |
| image_url | varchar(512) | | レースのヘッダー画像のURL (Cloudinary) |
| is_holding | boolean | NOT NULL | 開催されているか |

CREATE TABLE races (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), name varchar(32) NOT NULL, start_at timestamp NOT NULL, end_at timestamp NOT NULL, point_a char(16), point_b char(16), point_c char(16), athlete uuid[], memo text, image_url varchar(512), is_holding boolean);

## groups（グループ）
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | uuid | PRIMARY KEY DEFAULT uuid_generate_v4() | グループID |
| name | varchar(32) | NOT NULL | グループ名 |
| description | text | | 概要 |

CREATE TABLE groups (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), name varchar(32) NOT NULL, description text);

## tokens（トークン）
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| token | char(64) | NOT NULL PRIMARY KEY | トークン |
| permissions | varchar(64)[] | | 権限一覧 |
| user_id | uuid | | ユーザーID |
| description | text | | 概要 |

CREATE TABLE tokens (token char(64) NOT NULL PRIMARY KEY, permissions varchar(64)[], user_id uuid, description text);
