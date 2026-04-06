# nms_metrics 

## Заметки

- device_id в конфиге нужно использовать uuid
- На базовой станции сервисы лежат в папке /opt
- На dev сервисы лежат в папке /srv и /opt

## Основной принцип работы

- MetricsSender по интервалу вырезает метрики-outbox и шлет их в edgeCollector (с концами).
- MetricsCleaner по интервалу удаляет метрики из таблиц metrics, metrics-outbox (относительно времени).
- StatPuller по интервалу запрашивает метрики от eNodeB. Тот в свою очередь шлет запрос (stats, ncell) на
  eNodeBWS, для получения метрик. Далее StatPuller кладет их все (+ системные метрики) в таблицы
  (metrics, metrics-outbox, metrics-outbox-legacy).
- По web-запросу на ручку /metrics вырезаются метрики из таблицы metrics-outbox-legacy.

## Адреса

https://wiki.yandex.ru/homepage/instrukcii/adresa-xostov/adresa-servisov/

БС 105 (device_id 9b8cd954-936e-4916-8b64-a9f8c50233d6) - 10.10.10.159
БС 103 (device_id f2ecb3bf-2a4d-45f4-bd38-6c7cb63020cb) - 10.10.10.118