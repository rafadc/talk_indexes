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

# Explain plans

The basic mechanism to understand how your query will run

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

<!-- _class: lead -->

![bg left](./assets/table-of-contents.jpg)

# Single column indexes

<!-- Image attribution: https://www.flickr.com/photos/o_0/26278975918 -->

---

# BTrees

---

# Cardinality

- An index is more effective the more rows it can discard

---

# Costs of an index

- Storage space
- Insertion time

<!-- In mySQL we can check the performance schema to cleanup indexes -->


---
<!-- _class: lead -->

# Multiple column indexes

![bg right](./assets/mandelbrot.jpg)

---
<!-- _class: lead -->

# Rules of thumb

![bg left](./assets/rule-of-thumb.jpg)

---

# Put first the column that will remove the most values

---

# The range condition of the query should be the last one

---

# Like queries are range queries

---

# A LIKE %something query cannot use the index

---

# Range queries


| Name   | Surname     | Date of birth    |
| ------ | ----------- | ---------------- |
| Rafael | de Castro   | 14-04-1980       |


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

# Tricking your users

- Make values mandatory
- Put wise defaults

