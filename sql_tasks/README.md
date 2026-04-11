# SQL Tasks — GUAP University Database

## Schema

```
students  ──< grades >── subjects
                              ^
schedule >─────────────── subjects
```

## Tasks

| # | Query | QA Purpose |
|---|-------|------------|
| 1 | Students in group Z3420 | Verify enrollment |
| 2 | Student count by group | Cross-check after import |
| 3 | Average grade per semester | Compare with UI |
| 4 | Students with failing grades | Check debt notifications |
| 5 | Duplicate emails | Uniqueness check |
| 6 | Monday schedule for Z3420 | Compare schedule in app vs DB |
| 7 | Students with no grades | Data integrity after session start |
| 8 | Subjects with no exams | Data integrity |
| 9 | Top-5 students by GPA | Verify ranking feature |
| 10 | Exams by weekday | Load analytics |
| 11 | NULL in critical fields | Post-migration check |
| 12 | Delete test data | Teardown after autotests |

## Practice online

- [db-fiddle.com](https://www.db-fiddle.com)
- [sqliteonline.com](https://sqliteonline.com)
