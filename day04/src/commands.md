./client-candy -k "AA" -c 2 -m 50 --host candy.tld  --port 42003
- запуск клиента на домене (локально присвоила candy.tld localhost в файле etc/hosts)  
- в качестве параметров передаются хост, порт, тип конфет и их количество, сумма денег  

./server-candy --tls-ca cert/minica.pem --tls-certificate cert/candy.tld/cert.pem --tls-key cert/candy.tld/key.pem  --tls-host candy.tld --tls-port 8888  
--tls-ca ../../cert/minica.pem --tls-certificate ../../cert/candy.tld/cert.pem --tls-key ../../cert/candy.tld/key.pem  


- запуск сервера с добавлением сертификатов  

minica -domains candy.tld
- генерирование сертификатов  

openssl x509 -in cert/candy.tld/cert.pem -noout -modulus | openssl md5  
openssl rsa -in cert/candy.tld/key.pem -noout -modulus | openssl md5  
- проверка соответствия ключа и сертификата - должны возвращать одинаковые хеши  

curl --cert cert/candy.tld/cert.pem --key cert/candy.tld/key.pem --cacert cert/minica.pem -X POST https://candy.tld:8888/buy_candy -d '{"money": 50, "candyType": "AA", "candyCount": 2}' -H "Content-Type: application/json"


swagger generate cli -f api/swagger.yaml  
генерация клиента   
есть еще swagger generate client -f api/swagger.yaml которая генерирует более сложный код клиента  


swagger generate server -f api/swagger.yaml  
генерация сервера
