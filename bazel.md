Bazel Compilation Instructions
==============================

Building with bazel means you can place the directory anywhere in your system.
It also double acts as dependency control system with versioning.

---

Building
--------

Instructions are pretty similar to the ones outlined in README.md

1. Make sure you have [Bazel][1] installed
2. Initialize the schema located in `schemas/schema.sql`
3. Build and run the scraper:

   ```shell
   bazel run //:scraper
   ```

4. Build and run the server:

   ```shell
   bazel run //server
   ```

---

Adding Dependencies
-------------------

Make sure to add any external dependencies in `WORKSPACE`. Don't forget to
then add them to the appropate `BUILD` files where they are required.

---

Complaints
----------

Complain to [@Voytechnology][2] - If there are any questions regarding Bazel,
add `@Voytechnology` in the comments.

[1]: https://bazel.build
[2]: https://github.com/voytechnology