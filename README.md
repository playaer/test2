## Задача

- Написать сервис на Golang, который принимает массив URL-ов в теле, 
для данных URL он должен загрузить инф. о кол-во тегов на странице, код ответа, 
и все заголовки, пример ответа:

``` json
[
  {
    "url": "http://www.example.com/",
    "meta": {
      "status": 200,
      
      "headers": [
        {
          "content-type": "text\/html",
          "server": "nginx",
          "content-length": 605,
          "connection": "close",
          // ...
        }
      ]
    },
    "elemets": [
      {
        "tag-name": "html",
        "count": 1
      },
      {
        "tag-name": "head",
        "count": 1
      },
      // ...
    ]
  },
  // ...
]
```
- Сервис необходимо завернуть в Docker

## Решение/проверка

``` bash
$ git clone git@gitlab.com:playaer/gotask.git
$ cd gotask
$ docker-compose up -d
$ curl -X POST -d "[\"http://tut.by\",\"http://google.by\"]" http://localhost:8080
```