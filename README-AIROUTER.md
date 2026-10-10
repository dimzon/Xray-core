# Xray-core — форк для AI Router

Форк [XTLS/Xray-core](https://github.com/XTLS/Xray-core) на теге `v26.9.30`.
Нужен потому, что апстрим за последние годы вырезал поддержку неподписанных
и нешифрованных подключений, а ноды из подписок такие используют. Переписать
ноды нельзя, поэтому патчится ядро.

Основная ветка форка: `airouter/no-tls-verify`.

## Ветка `airouter/throne`

Те же патчи, но поверх пина `e662c22` из `throneproj/xray-core` — дерева,
которое линкует демон Throne (`core/go.mod` подменяет `xtls/xray-core` своим
форком). Черри-пик лёг без конфликтов: код совпадает. Их слой egress-интеграции
(`core/throne.go`, `transport/internet/throne.go` и др.) не тронут.

Чтобы собрать `ThroneCore.exe` с нашими патчами, в `throneproj/Throne`,
каталог `core`:

```
protoc -I gen --go_out=gen --go-grpc_out=gen gen/libcore.proto
```

затем в `core/go.mod` подменить `replace` на наш форк. Проще всего так —
номер версии Go посчитает сам:

```
go mod edit -replace github.com/xtls/xray-core=github.com/dimzon/Xray-core@airouter/throne
go mod tidy
```

и обычная сборка демона. Проверено: `ThroneCore.exe` (77 МБ) собирается
и стартует на этом дереве.

## Чем отличается от апстрима

| Изменение | Откат какого коммита | Зачем |
|---|---|---|
| `allowInsecure` в `tlsSettings` снова работает | `55956f8` (там же была таймбомба: флаг умирал после 2026-06-01) | Ноды с самоподписанными сертами |
| Убраны варнинги deprecated-фич (WS, VMess, Trojan, SS, gRPC) | — | Спам в логе на каждый запуск, ноды не переписать |
| VLESS/Trojan без шифрования на публичные адреса больше не роняют конфиг | `d7fa207` (2026-07-07) | Одна такая нода убивала загрузку всего конфига |
| VMess `security=none/zero` | `d7fa207` | Ноды `vmess://...?security=none` из подписок |
| Shadowsocks `method=none/plain` | `d7fa207` | Ноды `ss://...#...` с пустым шифром |
| Shadowsocks блочные шифры `aes-128-cfb`, `aes-256-cfb`, `chacha20`, `chacha20-ietf` | `00bcd40c` (2021-09-17) | Старые ноды из подписок |

Что НЕ тронуто: ошибки удалённых фич (`PrintRemovedFeatureError`), рефактор
`freedom.go` на `geodata`-матчеры, номера уже занятых enum-значений.
Возвращённые номера совпадают с историческими (`allow_insecure=1`,
`NONE=9`, `AES_128_CFB=1`…`CHACHA20_IETF=4`, `NONE=5`/`ZERO=6` у VMess),
поэтому со старыми бинарями и конфигами совместимо.

Дескрипторы `.pb.go` пересобраны программно через protobuf-go (протокольного
компилятора в системе нет): поле вставляется в дескриптор, структура и геттер
дописываются, затем `round-trip` маршалинга и существующие тесты пакета.
Подробности приёмов — в сообщениях коммитов ветки.

## Сборка локально

Нужен Go из `go.mod` (сейчас 1.27). Тег собирается так же, как апстрим:

```
go build -trimpath -buildvcs=false -o xray.exe -ldflags "-s -w -buildid=" ./main
```

`CGO_ENABLED=0` для релизных сборок (как в пайплайне).

Проверка без запуска:

```
xray.exe run -c config.json -test
```

Полезные тесты после правок:

```
go test ./transport/internet/tls/ ./proxy/shadowsocks/ ./proxy/vmess/... ./common/protocol/
go test ./testing/scenarios/ -run 'TestVMessNone|TestVMessZero|TestShadowsocksNone|TestShadowsocksStreamCipher'
```

`TestECHDial` требует прямого интернета и падает в закрытых сетях — это не
наше, он падает и на чистом дереве.

## Релизы

Пайплайн `.github/workflows/airouter.yml`:

- пуш в `airouter/**` — сборка `windows-amd64` + `linux-amd64`, артефакты;
- тег вида `v26.9.30-airouter1` — то же плюс GitHub Release с бинарниками;
- в версию бинаря подставляется короткий SHA с суффиксом `-airouter`
  (`xray.exe version` покажет `b26a91d-airouter` вместо чистого `b26a91d`).

## Обновление с апстрима

``` 
git fetch upstream
git checkout -b airouter/merge-<версия> <новый-тег>
git merge <новый-тег>   # конфликты — руками, наши правки помечены AIROUTER
```

Все правки форка помечены комментарием `// AIROUTER:`, искать по нему.
Файлы `.pb.go` при мерже лучше пересобирать заново тем же приёмом, а не
мержить построчно: там blob дескриптора.
