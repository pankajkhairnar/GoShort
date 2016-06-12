# goShort
GoShort is a URL shortener written in Golang and [BoltDB] for persistent key/value storage and for routing it's using high performent [HTTPRouter]

### REST Endpoints
  - [POST] - http://localhost:8080/create/ : It accepts POST form data with parameter "url" and returns json response with short URL
  - [GET] - http://localhost:8080/{SHORT_CODE}/ : If SHORT_CODE is valid and found in db request will be redirected to original URL
  - [GET] - http://localhost:8080/{SHORT_CODE}/json : If SHORT_CODE is valid and found in db it will return original url in json response  

### License
[Apache]


   [HTTPRouter]:<https://github.com/julienschmidt/httprouter>
   [BoltDB]: <https://github.com/boltdb/bolt>
   [Apache]: <https://github.com/pankajkhairnar/goShort/blob/master/LICENSE>
   
   
   
