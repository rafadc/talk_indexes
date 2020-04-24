---
theme: gaia
_class: lead
paginate: true
backgroundColor: #fff
backgroundImage: url('assets/bg.png')
---
<!-- _class: lead -->

# Database indexes

---

# mySQL

<!-- We will use mySQL as a mean to provide examples but this should be easily applicable to any other DB -->

---

# Your query is slow

- As your tables grow and grow, accessing your data will be slower and slower
---

# Explain plans

- The basic mechanism to understand how your query will run

---

# Full scans on big tables

- Normally we want to avoid
- On absurdly bigtables maybe you need different mechanisms outside the scope of this talk

---
![bg left](./assets/money_burn.jpg)

# Speeding up queries with indexes

---

![bg right](./assets/mark_twain.jpg)
# Single column indexes

---

# BTrees

---
<!-- _class: lead -->

# Multiple column indexes

---

# Ordering columns

- Put first the column that will remove the most values

---

# Range queries


| Name   | Surname     | Date of birth    |
| ------ | ----------- | ---------------- |
| Rafael | de Castro   | 14-04-1980       | 


---

# Turning range queries into specific value queries

---

# Tricking your users

- Make values mandatory
- Put wise defaults

