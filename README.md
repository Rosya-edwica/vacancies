# vacancies
Super Golang Parser 2000

## Описание
Эта программа создана для того, чтобы объединить в себе все парсеры вакансий. Идея заключалась в том, чтобы сделать общий проект, в котором будут прописаны общие настройки работы с БД, логгированием, уведомлениями и моделями. А парсеры будут как дополнения, которые будут реализовывать единый интерфейс API

## Запуск программы
Для того чтобы запустить программу, необходимо в консоли прописать:

```console
go run cmd/main.go headhunter area
```

1. `headhunter` можно заменить на `trudvsem` или `superjob`
2. `area` - необязательный параметр, используется для того, чтобы прочитать названия профобластей из файла prof_areas.txt

## P.S.
Одна из причин написать эту программу - изучение интерфейсов, поэтому не судите строго начинающего Golang-разработчика :)