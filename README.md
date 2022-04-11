# Сервер

В папке Server представлен код реализации gRPC прокси-сервиса для загрузки thumbnail’ов (превью видеоролика) c видеороликов YouTube. 
При повторном запросе на тот же видеоролик сервис отдает закэшированный ответ (в качестве хранилища используется sqlite). 

Сборка производилась компилятором Go 1.18

Запуск производится без параметров.
Принимает следующие флаги:

--addr        - Адрес сервера в формате host:port. 
                Значение по умолчанию: localhost:50051

--port        - Порт сервера. Значение по умолчанию: 50051.

--database    - Имя файла базы данных кэша. Значение по умолчанию: "".
                В случае значения по умолчанию файл базы данных не 
                будет создан. Кэш останется в оперативной памяти.


# Клиент

В папке Сlient представлен основной код реализации клиенклиентской части - утилита коммандной строки, которая принимает в качестве параметров ссылки на видеоролики и сохраняет полученные thumbnail’ы. Имена сохраненных thumbnail’ов являются (уникальной) частью url соответствующего видеоролика YouTube.

Сборка производилась компилятором Go 1.18

При запуске в качестве параметров принимает ссылки на видеоролики YouTube.

Принимает следующие флаги:

--async         - Файлы загружаются асинхронно, если true, иначе по порядку.

--out           - Папка загрузки thumbnail’ов. По умолчанию текущая папка
                  (откуда была запущена утилита).

--serverAddr    - Адрес сервера для получения RPC "GetThumbnail" в формате 
                  host:port. Значение по умолчанию: localhost:50051

Пример запуска в командной строке Windows:
>client --async --out=C:\Users\%USERNAME%\Downloads https://youtu.be/CWcj99dc650 https://youtu.be/_bL0s9JRVRk https://youtu.be/duy7bSyPLhs           


# Интерфейс взаимодействия 
Интерфейс взаимодействия клиента и сервера описан в файле thumbnail/thumbnail.proto. Код интерфейсов  gRPC "сервера" и "заглушки/клиента" сгенерированы protoc-3.20.0

# Автор

Никита Шеремета
guamoko95@gmail.com