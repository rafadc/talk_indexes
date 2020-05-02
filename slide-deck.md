---
theme: gaia
_class: lead
paginate: true
backgroundColor: #fff
backgroundImage: url('assets/bg.png')
---
<!-- _class: lead -->

# DB Indexing for performance

---

<!-- _class: lead -->

![10%](assets/mysql-logo.png)

<!-- We will use mySQL as a mean to provide examples but this should be easily applicable to any other DB -->

---

# We will oversimplify

<!-- In a lot of situations we will do simplifications. This is not an advanced talk -->


---

# Your query is slow

- As your tables grow and grow, accessing your data will be slower and slower

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

....

---

# Explain plans

- The basic mechanism to understand how your query will run

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

- *possible_keys*: Indexes that are applicable to your query
- *key*: Index that mySQL decided that is the best for this query
- *rows*: Number of rows in the index
- *extra*: More information on how the index works


---

# Full scans on big tables

- Normally we want to avoid
- On absurdly bigtables maybe you need different mechanisms outside the scope of this talk

---

# Single column indexes

---

# BTrees

---

# Cardinality

- An index is more effective the more rows it can discard

---

# Costs of an index

- Storage space
- Insertion time

---
<!-- _class: lead -->

# Multiple column indexes

---

# Rules of thumb

- Put first the column that will remove the most values
- The range condition of the query should be the last one
- Like indexes are range indexes
- A %something query cannot use the index

---

# Range queries


| Name   | Surname     | Date of birth    |
| ------ | ----------- | ---------------- |
| Rafael | de Castro   | 14-04-1980       |


---

# Turning range queries into specific value queries

- A classical trick is to turn a

``` sql
WHERE A > 1 AND A < 3
```

into

``` sql
WHERE A IN (1, 2, 3)
```


---

# Tricking your users

- Make values mandatory
- Put wise defaults

