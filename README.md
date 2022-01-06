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
| imei | char(15) | NOT NULL | IMEI |
| name | varchar(32) | NOT NULL | 端末名 |
| model | varchar(32) | NOT NULL | モデル番号 |
| lat | float | | 緯度 |
| lng | float | | 経度 |

## users (ユーザー情報)
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | char(16) | NOT NULL | ユーザーID |
| login_name | varchar(16) | NOT NULL | ログインID |
| display_name | varchar(32) | NOT NULL | 名前 |
| password | varchar(16) | NOT NULL | パスワード |
| group_id | char(16) | NOT NULL | グループID |
| user_type | varchar(16) | NOT NULL | `athlete`, `admin` or 'developer'
| device_imei | char(16) | | デバイスIMEI |
| sail_num | int(2) | | セイル番号 |
| course_limit | float | | コースリミット |
| image | varchar(512) | | プロフィール画像のURL (Cloudinary) |
| note | text | | 備考 |

## races (レース情報)
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | char(16) | NOT NULL | レースID |
| name | varchar(32) | NOT NULL | レース名 |
| start | datetime | NOT NULL | 開始時刻 |
| end | datetime | NOT NULL | 終了時刻 |
| point_a | char(16) | | A地点のデバイスIMEI |
| point_b | char(16) | | B地点のデバイスIMEI |
| point_c | char(16) | | C地点のデバイスIMEI |
| athlete | char(16)[] | | 競技者(ユーザー)ID一覧 |
| memo | text | | メモ |
| image | varchar(512) | | レースのヘッダー画像のURL (Cloudinary) |
| is_holding | boolean | NOT NULL | 開催されているか |

## groups（グループ）
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| id | char(16) | NOT NULL | グループID |
| name | varchar(32) | NOT NULL | グループ名 |
| description | text | | 概要 |

## tokens（トークン）
| Column | Type | Options | Description |
| ------ | ---- | ------- | ----------- |
| token | char(64) | UNIQUE | トークン |
| permissions | varchar(64)[] | | 権限一覧 |
| user_id | char(16) | | ユーザーID |
| description | text | | 概要 |

CREATE TABLE tokens (token char(64) UNIQUE, permissions varchar(64)[], user_id char(16), description text);
