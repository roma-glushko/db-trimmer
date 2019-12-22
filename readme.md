# DB-Trimmer

DB-Trimmer is a tool to reduce database size for staging or development environments in an intelligent way.

## PoC Benchmark

### Preconditions

MacBook Pro 15" 2017
Intel Core i7 (2.8Ghz)
16 GB memory

~913,660 Magento catalog database with simple, configurable and bundle products

### Plain Delete Query

```sql
DELETE FROM catalog_product_entity;
```
~913,660 - X

### Blocking 2k-chunks Deleting

~913,660 records - ~55mins

### Blocking 1k-chunks Deleting

~913,660 records - ~15mins

### Blocking 500-chunks Deleting

~913,660 records - 30mins

### Async (2 threads) 1k-chunks Deleting

~913,660 records - 28mins

### Async (4 threads: 2 for planning and 2 for trimming) 1k-chunks Deleting

~913,660 records

Planning - X
Trimming - X