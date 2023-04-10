## How to run project

- Import database that available in file support/migration.sql
- set mock in mockoon with endpoint 
    - api/v1/shopee/find-data/1234123412341 (get value from file support/mock_find_data_shopee.json)
    - api/v1/order (get value from file support/mock_store_data_shopee.json)
- change env SHOPEE_POST_ORDER and SHOPEE_FIND_DATA depends on your mockoon
- insert command make run