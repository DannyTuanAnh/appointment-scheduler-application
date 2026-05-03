# Appointment Scheduler Application

## Project Overview

**Appointment Scheduler Application** is a backend service for managing vehicle maintenance appointments at car dealerships. The system is designed to coordinate three core resources—**dealership working hours**, **service bays**, and **technicians**—to ensure appointments are scheduled efficiently and without conflicts.

### Key Capabilities

- **Dealership management**
  - Store dealership information and working hours (open/close time).
- **Service management**
  - Define services with expected duration and required bay type.
  - Attach required skills for each service (service requirements).
- **Technician management**
  - Manage technicians by dealership, skill sets, and active/inactive status.
  - Support operational actions such as transfer between dealerships and leave/return-to-work.
- **Service bay management**
  - Manage service bays by dealership and bay type with active/inactive status.

### Scheduling & Concurrency Safety (No Overlapping Appointments)

To support real-world concurrent bookings (multiple users creating appointments at the same time), the application enforces **conflict prevention at the database layer**:

- The `appointments` table uses PostgreSQL `tstzrange` (`duration`) to represent booking windows as **inclusive start, exclusive end**: `[start, end)`.
- GiST indexes accelerate time-range queries and availability checks.
- **Exclusion constraints** ensure that:
  - A **service bay** cannot be booked for overlapping time ranges.
  - A **technician** cannot be booked for overlapping time ranges.
  - Cancelled appointments are ignored by conflict constraints so the time slot becomes available again.

This approach guarantees correctness even under high concurrency, because the database becomes the final authority that rejects overlapping inserts/updates.

### API Design

Routes are intentionally split into two categories:

- **App routes (`/api/v1/app/...`)**: user-facing endpoints for reading/searching and booking.
- **Admin routes (`/api/v1/admin/...`)**: administrative endpoints for create/update/delete operations.

### Architecture

The project follows a Clean Architecture-inspired structure:

**Route → Handler (Gin) → Service → Repository → sqlc (PostgreSQL)**

- Handlers manage HTTP binding/validation and produce consistent responses.
- Services implement business logic and cross-table validations.
- Repositories encapsulate database access and translate DB errors into application errors.
- `sqlc` generates type-safe queries for PostgreSQL.

## How To Run

### Option A: Using Docker and DBeaver (recommended)

1. Clone this repository.
2. Start PostgreSQL using Docker Compose:

   ```bash
   docker-compose up -d
   ```

3. Use DBeaver to connect to the database using the credentials in `docker-compose.yaml`.
4. Run the API server (see **Run the application** below).

### Option B: Using local PostgreSQL (PgAdmin)

1. Clone this repository.
2. Create a PostgreSQL database.
3. Import the schema SQL file (e.g. `schema.sql`) into your database using PgAdmin.
4. Update environment variables in `.env`.
5. Run the API server (see **Run the application** below).

### Run the application

- **Run from source:**

  ```bash
  go run ./cmd/api
  ```

- **Run the prebuilt executable:**
  1. Ensure `.env` is placed next to `main.exe` (or that your system environment variables are set).
  2. Double-click `main.exe` to start the server.

## How to Test

> Prerequisites:
>
> - The API server is running on `http://localhost:8080`
> - Database is migrated and seeded with sample data (dealership/service/bay/technician)
> - Time zone in examples uses `+07:00`

### 1) Check Current Availability

Use the availability endpoint to retrieve the dealership’s current working window and already occupied time ranges for the selected service.
This endpoint is used by the client to render a booking timeline and guide the user toward valid appointment times.

```bash
  curl --location --request GET 'http://localhost:8080/api/v1/app/appointment?dealership_id=1&service_id=2&preference_time=2026-05-02&bay_type_id=2&anticipated_duration=120'
```

Expected Response (example)

```json
{
  "data": {
    "date": "02-05-2026",
    "work_start": "08:00",
    "work_end": "18:00",
    "duration_minutes": 120,
    "busy": [
      {
        "start": "11:00",
        "end": "15:00"
      }
    ]
  },
  "status": "success"
}
```

What to verify

- the dealership operating window is returned correctly
- the requested service duration is reflected
- already occupied time ranges are returned
- the frontend can derive valid time slots from this response

### 2) Create a Valid Appointment

Submit a valid booking request using a time range that does not overlap with an existing occupied period.

```bash
curl --location --request POST 'http://localhost:8080/api/v1/app/appointment' \
--header 'Content-Type: application/json' \
--data '{
  "dealership_id": 1,
  "service_id": 2,
  "bay_type_id": 2,
  "customer_name": "test",
  "start_time": "2026-05-02T08:00:00+07:00",
  "end_time": "2026-05-02T11:00:00+07:00"
}'
```

Expected Result

- booking is accepted
- appointment is persisted
- technician and service bay are assigned
- appointment becomes visible in subsequent availability queries

What to verify

- response returns success
- database contains a new appointment
- the booked range is now marked as occupied

### 3) Verify Availability Changes After Booking

Call the availability endpoint again after a successful booking.

Expected Result

The previously booked time range should now appear in the `busy` list.

This confirms that:

- the appointment was persisted
- availability is recalculated correctly
- subsequent users will no longer see the slot as free

### 4) Verify Double-Booking Protection

Attempt to create another appointment for the same dealership and overlapping time range.

```bash
curl --location --request POST 'http://localhost:8080/api/v1/app/appointment' \
--header 'Content-Type: application/json' \
--data '{
  "dealership_id": 1,
  "service_id": 2,
  "bay_type_id": 2,
  "customer_name": "duplicate-test",
  "start_time": "2026-05-02T09:00:00+07:00",
  "end_time": "2026-05-02T10:30:00+07:00"
}'
```

Expected Result

- request is rejected with HTTP `409 Conflict`
- no duplicate appointment is created
- overlapping resources are not assigned twice

What to verify

- API returns conflict / validation failure
- no duplicate appointment is created
- overlapping resources are not assigned twice

### 5) Run Concurrency Test

Simulate multiple clients attempting to book the same slot concurrently.

Goal

Validate that the system prevents race-condition double booking under concurrent load.

Expected Result

- only one request succeeds
- all other requests fail gracefully
- only one appointment record is persisted for the slot

Run

```bash
go test ./tests/concurrency -v
```

What to verify

- no duplicate slot allocation occurs
- database remains consistent
- contention is handled safely

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
