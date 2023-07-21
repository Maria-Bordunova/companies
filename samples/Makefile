MAKEFLAGS += --always-make

SERVICE_HOST:="localhost:7150"

create-company:
	curl -X POST -H "Content-Type: application/json" \
	-H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2ODk5NDUxMjgsImV4cCI6MTcyMTQ4MTEyOCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.KcdMPrS69yZ05E4vfPq_dtM1fSiuQKWYlqrOtqEJMys" \
	-d '{ "id": "12345", "name": "Sample Company 2","description": "A sample company for testing","employees": 100,"registered": true,"type": "corporations"}' \
	http://localhost:7150/companies

get-company:
	curl -X GET http://localhost:7150/companies/12345

update-company:
	curl -X PATCH -H "Content-Type: application/json" \
	-H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2ODk5NDUxMjgsImV4cCI6MTcyMTQ4MTEyOCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.KcdMPrS69yZ05E4vfPq_dtM1fSiuQKWYlqrOtqEJMys" \
	-d '{"name": "Updated Company Name", "employees": 150 }' \
    http://localhost:7150/companies/12345

delete-company:
	curl -X DELETE \
	-H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2ODk5NDUxMjgsImV4cCI6MTcyMTQ4MTEyOCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.KcdMPrS69yZ05E4vfPq_dtM1fSiuQKWYlqrOtqEJMys" \
	 http://localhost:7150/companies/12345
