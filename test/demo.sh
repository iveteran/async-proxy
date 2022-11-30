curl -v -H "app-id:fimatrix" -H "app-token:fimatrix2020" -H "x-uid: 9" \
  -H "proxy-backend:http://localhost:8090/" \
  -H "proxy-topic:fap_demo" \
  -d '{"name": "yufangbin", "age": 20}' \
http://localhost:4035/dp/fdtx/extract_file_tables

curl -v -H "app-id:fimatrix" -H "app-token:fimatrix2020" -H "x-uid: 9" \
  -d '{"name": "yufangbin", "age": 20}' \
http://localhost:4035/dp/fdtx/extract_file_tables
