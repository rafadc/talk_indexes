---
theme: gaia
class:
  - lead
  - invert
paginate: false
---
<!-- _class: lead -->

# DB Indexing for performance

---

<!-- _class: lead -->

![bg](assets/mysql-logo.png)

<!-- We will use mySQL as a mean to provide examples but this should be easily applicable to any other DB -->

---

# We will oversimplify

![bg left](assets/kid-reading.jpg)

<!-- In a lot of situations we will do simplifications. This is not an advanced talk -->

---

# Batteries included

You have a docker project that prepares the environment of this talk

[https://github.com/rafadc/talk_indexes](https://github.com/rafadc/talk_indexes)

Feel free to use it as a playground to experiment

---

# Your queries become slow

```sql
SELECT * FROM people_small WHERE name="John" AND company = "Mertz-Mertz";

....

took 7ms
```


```sql
SELECT * FROM people_without_indexes WHERE name="John" AND company = "Mertz-Mertz";

....

took 3.8s
```

<!-- Normally queries don't start slow. They gradually degrade through time -->

---

<!-- _class: lead -->

![bg left](./assets/indexes.jpg)

# Indexes

<!-- Image attribution: https://www.flickr.com/photos/gotcredit/33756630285/in/photolist-TqXy7Z-41UtU-65pRH6-74qEZ5-74qFbQ-74qERU-74qEHb-74mLdZ-74mL4v-4rcRjS-74qF5N-Mcgq6s-25XMDFQ-29iao55-29jdUEd-e85NKw-29nixZp-4VJpnq-e85P8w-bfVCgn-bfVwrr-bfVypK-6wwm7H-66kvGn-KKAUxv-PNbZEk-LVgyVT-ejfuDg-oi246q-bfVogR-bfVqsH-bfVskt-bfVub2-bfVAyT-wLoYGz-4GX4gt-G3ANZh-omNTSx-ok1QQ7-o3yiZK-ojRrsL-o3yhNX-ok1Qzh-omNSkV-o3y2Wd-ojLeGZ-o3ykPA-L4gQTo-LemtJ2-od3DLC -->

---

# Index

A DB structure to retrieve values more quickly

---

# Index

As a rule of thumb one query uses one index

<!-- This is not true at all. We have union indexes or merge indexes but we will simplify our mental model -->

---

# Force index

```sql
SELECT *
FROM people_multi_column_index
USE INDEX (people_multi_column_index_happy_name_IDX)
WHERE name = "john" AND happy = true;
```

<!-- You can always manually choose an index but it ismore interesting to know when the query optimizer pìcks an index -->

---

# ANALYZE TABLE

Performs a key distribution analysis and stores the distribution for the named table or tables

```sql
ANALYZE TABLE people_single_index;
```

Remember to run when changing indexes in your playground

<!-- So remember to analyze table when creating indexes in your playground -->

---

# Explain plans

The basic mechanism to understand how the query optimizer will run your query

```sql
EXPLAIN SELECT * FROM people_without_indexes WHERE name="John" AND company = "Mertz-Mertz";
```

```shell
id|select_type|table                 |partitions|type|possible_keys|key|key_len|ref|rows   |filtered|Extra      |
--|-----------|----------------------|----------|----|-------------|---|-------|---|-------|--------|-----------|
 1|SIMPLE     |people_without_indexes|          |ALL |             |   |       |   |9417967|       1|Using where|
```

---

# Reading an explain plan

```shell
id|select_type|table                 |partitions|type|possible_keys|key|key_len|ref|rows   |filtered|Extra      |
--|-----------|----------------------|----------|----|-------------|---|-------|---|-------|--------|-----------|
 1|SIMPLE     |people_without_indexes|          |ALL |             |   |       |   |9417967|       1|Using where|
```

- *type*: Type of join to access the table
- *possible_keys*: Indexes that are applicable to your query
- *key*: Index that mySQL decided that is the best for this query
- *rows*: Estimation of rows to be examined
- *filtered*: Percentage of rows filtered by conditions
- *extra*: More information on how the index works

---

<!-- _class: lead -->

![bg left](./assets/table-of-contents.jpg)

# Single column indexes

<!-- Image attribution: https://www.flickr.com/photos/o_0/26278975918 -->

---

# BTrees

![fit](./assets/btree.jpg)

---

# Cardinality

Uniqueness of values stored in a specified column within an index
An index is more effective the more rows it can discard

<!-- The primary key has cardinality equal to the number of rows. That makes it the most effective index to access individually -->

---

# Costs of maintaining an index

- Storage space
- Insertion time

<!-- In mySQL we can check the performance schema to cleanup indexes -->

---

# A simple query

```sql
SELECT * FROM people_without_indexes WHERE name="John" AND company = "Mertz-Mertz";

....

took 3.8s
```

```sql
CREATE INDEX people_single_index_name_IDX USING BTREE ON indexes.people_single_index (name);
```


```sql
SELECT * FROM people_single_index WHERE name="John" AND company = "Mertz-Mertz";

....

took 124ms
```


<!-- In the most simple scenarios an index can bring your speed back -->


---
<!-- _class: lead -->

# Multiple column indexes

![bg right](./assets/mandelbrot.jpg)

---


# Multiple column indexes

Sometimes we have not enough with filtering in only one column
The order of the index fields is paramount

---

# Multi column BTrees

![fit](./assets/btree.jpg)


---

# Candidate indexes

In order to use a column in the index its previous column needs to be in the WHERE clause

---

# Candidate indexes

```sql
CREATE INDEX people_multi_column_index_name_company_IDX USING BTREE ON indexes.people_multi_column_index (name, company);
```

```sql
EXPLAIN SELECT * FROM people_multi_column_index WHERE company = "Mitchell PLC";
```

```
id|select_type|table                    |type|possible_keys|key|key_len|ref|rows   |filtered|Extra      |
--|-----------|-------------------------|----|-------------|---|-------|---|-------|--------|-----------|
 1|SIMPLE     |people_multi_column_index|ALL |             |   |       |   |9700871|      10|Using where|
```



---

# Candidate indexes

```sql
CREATE INDEX people_multi_column_index_name_company_IDX USING BTREE ON indexes.people_multi_column_index (name, company);
```

```sql
EXPLAIN SELECT * FROM people_multi_column_index WHERE name = "John" AND company = "Mitchell PLC";
```

```
id|select_type|table                    |type|possible_keys                                    |key                                       |key_len|ref        |rows|filtered|Extra|
--|-----------|-------------------------|----|-------------------------------------------------|------------------------------------------|-------|-----------|----|--------|-----|
 1|SIMPLE     |people_multi_column_index|ref |people_multi_column_index_name_happy_IDX         |people_multi_column_index_name_company_IDX|204    |const,const|   1|     100|     |
                                              people_multi_column_index_name_date_of_birth_IDX
                                              people_multi_column_index_name_company_IDX
```



---

# The distribution of values is important

```sql
SELECT COUNT(*) FROM people_multi_column_index WHERE happy = true;
```

```
9998964
```

---

# The distribution of values is important

Slow

```sql
SELECT *
FROM people_multi_column_index
WHERE company = "Sample" AND happy = true;

took 55s
```

---

# The distribution of values is important

Fast

```sql
SELECT *
FROM people_multi_column_index
WHERE company = "Sample" AND happy = false;

took 10ms
```

☣ There is no indication in the explain plan ☣

<!-- The query optimizer has no way of knowing this -->

---

# The range condition of the query should be the last one in the index

```sql
SELECT *
FROM people_multi_column_index
USE INDEX (people_multi_column_index_date_of_birth_name_IDX)
WHERE name = "john" AND date_of_birth > "2000-01-01";

took 1s.
```

<!-- If you think on the BTree structure we cannot use a range and then another value -->

---

# The range condition of the query should be the last one in the index

```sql
SELECT *
FROM people_multi_column_index
WHERE name = "john" AND date_of_birth > "2000-01-01";

took 3ms.
```

☣ There is no indication in the explain plan ☣

---

# Like queries are range queries

Same rules as with a regular range query apply

```sql
SELECT *
FROM people_multi_column_index
WHERE name LIKE "john%" AND date_of_birth > "2000-01-01";
```

---

# A "LIKE %something" query cannot use a BTree

```sql
EXPLAIN SELECT * FROM people_single_index WHERE name LIKE "%John";
```

```
id|select_type|table              |partitions|type|possible_keys|key|key_len|ref|rows   |filtered|Extra      |
--|-----------|-------------------|----------|----|-------------|---|-------|---|-------|--------|-----------|
 1|SIMPLE     |people_single_index|          |ALL |             |   |       |   |9676556|   11.11|Using where|
 ```

---

# A "LIKE %something" query cannot use a BTree

```sql
CREATE FULLTEXT INDEX people_full_text_search_name_IDX ON indexes.people_full_text_search (name) WITH PARSER ngram;

EXPLAIN SELECT * FROM people_full_text_search WHERE MATCH(name) AGAINST ("ohn");
```

```
id|select_type|table                  |partitions|type    |possible_keys                   |key                             |key_len|ref  |rows|filtered|Extra                        |
--|-----------|-----------------------|----------|--------|--------------------------------|--------------------------------|-------|-----|----|--------|-----------------------------|
 1|SIMPLE     |people_full_text_search|          |fulltext|people_full_text_search_name_IDX|people_full_text_search_name_IDX|0      |const|   1|     100|Using where; Ft_hints: sorted|
 ```

---

# Turning range queries into specific value queries

``` sql
SELECT * WHERE A > 1 AND A < 3
```

into

``` sql
SELECT * WHERE A IN (1, 2, 3)
```


---

# Index only queries

Sometimes you don't even need to go to the table

```sql
SELECT *
FROM people_multi_column_index
WHERE name = "john";
```

```
id|select_type|table                    |partitions|type|possible_keys |key                                     |key_len|ref  |rows|filtered|Extra|
--|-----------|-------------------------|----------|----|--------------|----------------------------------------|-------|-----|----|--------|-----|
 1|SIMPLE     |people_multi_column_index|          |ref |    [...]     |people_multi_column_index_name_happy_IDX|102    |const|3540|     100|     |
```

---

# Index only queries

Sometimes you don't even need to go to the table

```sql
SELECT name, date_of_birth
FROM people_multi_column_index
WHERE name = "john";
```

```
id|select_type|table                    |partitions|type|possible_keys |key                                             |key_len|ref  |rows|filtered|Extra      |
--|-----------|-------------------------|----------|----|--------------|------------------------------------------------|-------|-----|----|--------|-----------|
 1|SIMPLE     |people_multi_column_index|          |ref |    [...]     |people_multi_column_index_name_date_of_birth_IDX|102    |const|3540|     100|Using index|
```


---

# Tricking your users to use your indexes more frequently

- Make values mandatory
- Put wise defaults

---

<!-- _class: lead -->

# Sorting

![bg right](./assets/order.jpg)

---

# A sort can use an index

```sql
EXPLAIN
SELECT *
FROM people_multi_column_index
WHERE name = "john" ORDER BY date_of_birth ASC;
```

```
id|select_type|table                    |partitions|type|possible_keys |key    |key_len|ref  |rows|filtered|Extra                |
--|-----------|-------------------------|----------|----|--------------|-------|-------|-----|----|--------|---------------------|
 1|SIMPLE     |people_multi_column_index|          |ref |  [...]       | [...] |102    |const|3540|     100|Using index condition|
```

---

# A sort can use an index

```sql
EXPLAIN
SELECT *
FROM people_multi_column_index
WHERE name = "john" ORDER BY date_of_birth DESC;
```

```
id|select_type|table                    |partitions|type|possible_keys |key    |key_len|ref  |rows|filtered|Extra      |
--|-----------|-------------------------|----------|----|--------------|-------|-------|-----|----|--------|-----------|
 1|SIMPLE     |people_multi_column_index|          |ref |  [...]       | [...] |102    |const|3540|     100|Using where|
```

⚠ Only if the index is using the same order ⚠

---


# More info

![bg left](assets/mysql-book.png)

