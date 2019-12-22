#!/usr/bin/env bash

##
# This script exports sample database information into a separate dumps
##

if [ "$1" == "-h" ] ; then
    echo "Usage: `basename $0` [-h] username password database [--timestamp]"
    exit 0
fi

USERNAME=$1
PASSWORD=$2
DATABASE=$3

##
# Building Query like
# SET group_concat_max_len = 10240; SELECT GROUP_CONCAT(table_name separator ' ') FROM information_schema.tables WHERE table_schema='DB' AND X
##

SQL_FETCH_TABLES="SET group_concat_max_len = 10240;"
SQL_FETCH_TABLES="${SQL_FETCH_TABLES} SELECT GROUP_CONCAT(table_name separator ' ')"
SQL_FETCH_TABLES="${SQL_FETCH_TABLES} FROM information_schema.tables WHERE table_schema='${DATABASE}'"

# Catalog
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_TABLES} AND ("

SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} table_name LIKE 'eav_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'catalog_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'sequence_product%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'sequence_catalog_category'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'url_rewrite'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'inventory_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'downloadable_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'layout_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'magento_catalogevent_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'magento_catalogpermissions_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'catalogrule_%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'magento_targetrule%'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'magento_targetrule%'"

SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'catalog_product_flat_cl'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'catalog_product_price_cl'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'atwix_catalogrule_product_cl'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'targetrule_product_rule_cl'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'catalogsearch_fulltext_cl'"
SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'customer_quote_versions_total_price_cl'"

SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES} OR table_name LIKE 'store%'"

SQL_FETCH_CATALOG_TABLES="${SQL_FETCH_CATALOG_TABLES})"

# Catalog
CATALOG_TABLE_LIST=`mysql -u${USERNAME} -p${PASSWORD} -AN -e"${SQL_FETCH_CATALOG_TABLES}"`

TIMESTAMP=""

if [[ $* == *--timestamp* ]] ; then
    TIMESTAMP=".$(date +%F-%T)"
fi

# Catalog
mysqldump --single-transaction --complete-insert --triggers -u${USERNAME} -p${PASSWORD} ${DATABASE} ${CATALOG_TABLE_LIST} | gzip > catalog-dump${TIMESTAMP}.sql.gz