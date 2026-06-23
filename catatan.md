#login dan register

curl -X POST http://localhost:8080/api/auth/register \
-H "Content-Type: application/json" \
-d '{"username": "fauzan", "email": "fauzan@leetcode.com", "password": "rahasia123"}'


curl -X POST http://localhost:8080/api/auth/login \
-H "Content-Type: application/json" \
-d '{"email": "fauzan@leetcode.com", "password": "rahasia123"}'

nyalakan ngrok untuk jenkins biar mengaktifkan webhook di github
./ngrok http 8081
