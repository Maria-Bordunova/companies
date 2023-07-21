MAKEFLAGS += --always-make

SERVICE_HOST:="localhost:7150"

create-company:
	curl -X POST -H "Content-Type: application/json" \
	-d '{ "id": "12345", "name": "Sample Company 2","description": "A sample company for testing","employees": 100,"registered": true,"type": "corporations"}' \
	http://localhost:7150/companies

get-company:
	curl -X GET http://localhost:7150/companies/12345

update-company:
	curl -X PATCH -H "Content-Type: application/json" \
	-d '{"name": "Updated Company Name", "employees": 150 }' \
    http://localhost:7150/companies/12345

delete-company:
	curl -X DELETE  http://localhost:7150/companies/12345
