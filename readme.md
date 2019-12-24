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

### Async (6 threads: 2 for planning and 4 for trimming) 1k-chunks Deleting

~913,660 records - ~40mins (~3secs avg delete query exec time)

### Async (6 threads: 2 for planning and 4 for trimming) 500-chunks Deleting

~913,660 records - 14mins (~1.8secs avg delete query exec time)

### Async (6 threads: 2 for planning and 4 for trimming) 300-chunks Deleting

~913,660 records - X (~2.5secs avg delete query exec time)

### Async (8 threads: 2 for planning and 6 for trimming) 300-chunks Deleting

~913,660 records - 8mins (~1.3secs avg delete query exec time)