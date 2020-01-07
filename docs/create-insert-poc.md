# Create & Insert PoC

Helpful MySQL commands to inspect database structure:
```sql
SHOW CREATE TABLE catalog_product_entity;
SHOW FULL TABLES;
SHOW COLUMNS FROM catalog_product_entity;
DESCRIBE catalog_product_entity;
SHOW INDEX FROM catalog_product_entity;
SHOW TABLE STATUS FROM cabinets_mage2 LIKE 'catalog_product_entity';

SHOW TRIGGERS LIKE 'catalog_product_entity';

SELECT TABLE_NAME, COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = 'cabinets_mage2' AND TABLE_NAME = 'catalog_product_entity'
```

go run main.go --db-pass "root" --db-name "db-trimmer-sample"
