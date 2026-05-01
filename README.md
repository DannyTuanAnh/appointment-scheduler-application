# Appointment Scheduler Application

## AI Collaboration Narrative

I utilized GenAI to assist in code generation. While the actual prompts were issued in my native language for precision, they have been translated below for clarity.

### Schema Design & Generation

#### Step 1: Initial DDL generation

I initially designed the system schema using an Entity Relationship Diagram in draw.io. Then, I provided this visual context to an AI coding tool with the following prompt:

- Generate SQL DDL scripts to create tables based on the provided design. Place each table's creation script in a corresponding `*.up.sql` file and the `DROP` statements in respective `*.down.sql` files.
- Implement indexes for query optimization, specifically for the `appointments` table which will be frequently queried for availability using `start_time`, `anticipated_time`, `bay_id`, and `technician_id`.
- Apply `UNIQUE` constraints to the `technician_skills` and `service_requirements` tables to prevent duplicate records.

#### Step 2: Validation & refactoring

Upon reviewing the AI's output, I identified an opportunity to optimize the appointment table. Instead of using separate `start_time` and `end_time` columns, I decided to use the PostgreSQL `tstzrange` data type combined with a GiST (Generalized Search Tree) index.

This approach is superior for:

- **Search performance**: Efficiently filtering out irrelevant time ranges.
- **Concurrency control**: Preventing overlapping schedules directly at the database level using exclusion constraints, while accounting for appointment statuses like `cancelled` or `no_show`.

The prompt I used for refactoring:

- Instead of using separate `start_time` and `anticipated_time` columns for appointment booking, could we utilize a single `duration` column with the `tstzrange` data type? The stored values would look like this: `[2023-10-01 08:00, 2023-10-01 09:00)`.
- Furthermore, we should implement a GiST (Generalized Search Tree) index to accelerate search performance by efficiently filtering out irrelevant time ranges.
- Most importantly, this approach enables overlap prevention directly at the database level, similar to how `UNIQUE` constraints work for IDs. However, please note that the overlap exclusion must be conditional based on the appointment status. Specifically, we need to prevent overlapping schedules unless the status is `cancelled` (where the slot is freed by the client) and `no_show`.

### SQL Query Implementation

First, I provided the table structures to the AI to establish context and requested the generation of CRUD operations.

Prompt (in English):

Please analyze the table structure provided below and execute the following requirements:

```sql
CREATE TABLE IF NOT EXISTS dealerships (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  open_time TIME NOT NULL,
  close_time TIME NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_dealerships_name ON dealerships (name);
```

- Generate queries to perform the following operations:
  1. Fetch all dealerships.
  2. Update all columns (except `id`) using `COALESCE` with `sqlc.narg(...)` filtered by `id`.
  3. Search for dealership information by name.
  4. Delete a dealership by its id.
  5. Other operations needed for business logic.

I carefully validated and manually adjusted the generated queries to ensure they aligned with the specific business logic of each module.

### Go (Golang)

To ensure the AI understood the project's Clean Architecture, I structured the core organization and the execution flow for the Dealership module as a reference. Then, I provided that context to the AI coding tool to learn it. Prompt:

- Please review and learn the code organization within these files to understand the execution flow, response handling, and DTO (Data Transfer Object) patterns. Pay close attention to how database validation is performed and how existing—as well as upcoming—utility functions are utilized.
- Note that the routes are clearly divided into two categories: user-facing routes and administrative routes. This understanding is essential for you to fulfill my subsequent requirements accurately.

By providing this high-level "Source of Truth," I was able to direct the AI to implement subsequent modules (like `skill_module.go`) while maintaining strict architectural consistency.

After that, I provided a SQL query file as context for the AI coding tool with the prompt:

> Proceed to implement the 6 corresponding routes for these queries in `skill_module.go`.
