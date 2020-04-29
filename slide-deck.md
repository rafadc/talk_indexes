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

# mySQL

<!-- We will use mySQL as a mean to provide examples but this should be easily applicable to any other DB -->

---

# We will oversimplify

---

# Your query is slow

- As your tables grow and grow, accessing your data will be slower and slower

---

# Explain plans

- The basic mechanism to understand how your query will run

```sql
EXPLAIN PLAN SELECT * FROM people_without_indexes WHERE company='Stuart';
```

---

# Reading an explain plan


---

# Full scans on big tables

- Normally we want to avoid
- On absurdly bigtables maybe you need different mechanisms outside the scope of this talk

---

# Speeding up queries with indexes

---

# Single column indexes

---

# BTrees

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

