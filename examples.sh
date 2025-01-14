curl -X POST http://localhost:8080/accounts \
-H "Content-Type: application/json" \
-d '{
  "name": "Electricity Bill",
  "amount": 150.50,
  "account_type": "payable"
}'

curl -X GET http://localhost:8080/accounts

curl -X PUT http://localhost:8080/accounts/1 \
-H "Content-Type: application/json" \
-d '{
  "name": "Water Bill",
  "amount": 75.25,
  "account_type": "payable"
}'

curl -X DELETE http://localhost:8080/accounts/1

