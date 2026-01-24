# Database Schema Documentation

## Overview

Am.EventKit uses PostgreSQL as its primary database with UUID-based primary keys and comprehensive role-based access control. The database supports a multi-role event management system with student registrations, team management, accommodation booking, and comprehensive payment tracking.

## Database Extensions

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";    -- UUID generation
CREATE EXTENSION IF NOT EXISTS citext;         -- Case-insensitive text
CREATE EXTENSION IF NOT EXISTS pg_cron;        -- Scheduled tasks
CREATE EXTENSION IF NOT EXISTS pg_duckdb;      -- Analytics capabilities
```

## Core Entity Types

### User Roles
- **Students**: Primary users who register for events and book accommodations
- **Organizers**: Departments/Clubs that create and manage events
- **Admin**: System administrators with full access
- **Accommodation Personnel**: Hospitality staff managing hostels and check-ins

### Event Types
- **EVENT**: Regular events
- **WORKSHOP**: Educational workshops

### Event Status
- **CLOSED**: No longer accepting registrations
- **ACTIVE**: Currently accepting registrations
- **COMPLETED**: Event has finished

## Table Structure

### User Management Tables

#### `student`
Primary user table for student registrations and authentication.

```sql
Columns:
- id: UUID (PK) - Unique identifier
- name: TEXT - Student's full name
- email: TEXT (UNIQUE) - Login email
- password: TEXT - Hashed password
- phone_number: TEXT - Contact number
- is_amrita_student: BOOLEAN - Internal student flag
- amrita_roll_number: TEXT - University roll number (nullable)
- college_name: TEXT - Default: 'Amrita Vishwa Vidyapeetham'
- college_city: TEXT - Default: 'Coimbatore'
- account_status: ENUM(VERIFIED, DISABLED) - Account state
- refresh_token: TEXT - JWT refresh token
- hospitality_id: TEXT (UNIQUE) - Accommodation ID
- created_at/updated_at: TIMESTAMP - Audit fields
```

**Design Rationale**: 
- Separate `hospitality_id` allows accommodation systems to operate independently
- College defaults reflect the primary institution while allowing external participants
- Roll number is nullable for non-university students

#### `student_onboarding`
Temporary table for OTP-based registration verification.

```sql
Columns:
- id: SERIAL (PK)
- name, email, password, phone_number: Basic info
- is_amrita_student, amrita_roll_number: University affiliation
- otp: TEXT - One-time password
- created_at/expiry_at: TIMESTAMP - OTP validity window
```

**Design Rationale**: Isolates registration process from main user table for security and cleanup purposes.

#### `password_reset`
Handles password recovery through OTP verification.

#### `admin`
System administrators with full platform access.

#### `organizer`
Event creators (departments/clubs) with hierarchical leadership structure.

```sql
Columns:
- id: UUID (PK)
- name: TEXT (UNIQUE) - Organization name
- email: TEXT (UNIQUE) - Login email
- org_type: ENUM(DEPARTMENT, CLUB) - Organization classification
- student_head: TEXT - Lead student organizer
- student_co_head: TEXT - Assistant student organizer
- faculty_head: TEXT - Faculty supervisor
```

**Design Rationale**: Supports university governance model with student-faculty leadership structure.

### Event Management Tables

#### `event`
Core event information with comprehensive metadata.

```sql
Columns:
- id: UUID (PK)
- name: TEXT (UNIQUE) - Event title
- blurb: TEXT - Short description
- description: TEXT - Full description
- cover_image_url: TEXT - Event image
- price: INTEGER - Registration cost in rupees
- is_per_head: BOOLEAN - Pricing model (flat vs per-person)
- rules: TEXT - Event rules/regulations
- event_type: ENUM(EVENT, WORKSHOP) - Category
- is_group: BOOLEAN - Team-based event
- is_technical: BOOLEAN - Technical event flag
- max_teamsize/min_teamsize: INTEGER - Team size limits
- total_seats: INTEGER - Capacity limit
- seats_filled: INTEGER - Current registrations
- event_status: ENUM(CLOSED, ACTIVE, COMPLETED)
- event_mode: ENUM(ONLINE, OFFLINE) - Delivery mode
- attendance_mode: ENUM(SOLO, DUO) - Attendance tracking
```

**Design Rationale**:
- Flexible pricing supports both flat-fee and per-person billing
- Team size constraints enable various event formats
- Technical flag helps categorize and filter events

#### `event_schedule`
Multiple schedule entries for multi-day events.

```sql
Columns:
- id: UUID (PK)
- event_id: UUID (FK) - References event
- event_date: DATE - Specific day
- start_time/end_time: TIMESTAMP - Event timing
- venue: TEXT - Physical/virtual location
```

#### `organizer` & `event_to_organizer_mapping`
Many-to-many relationship allowing multiple organizers per event.

**Design Rationale**: Enables collaboration between departments and clubs.

#### `tags` & `event_tag_mapping`
Flexible categorization system with abbreviations for UI display.

#### `people` & `people_to_event_mapping`
Guest speakers, judges, and special participants with day-specific availability.

### Registration and Booking Tables

#### `bookings`
Transaction and registration management.

```sql
Columns:
- id: UUID (PK)
- txn_id: TEXT (UNIQUE) - Payment gateway transaction ID
- student_id: UUID (FK) - Registered student
- event_id: UUID (FK) - Target event
- registration_fee: INTEGER - Amount paid
- registration_fee_without_gst: INTEGER - Base amount
- product_info: TEXT - Payment description
- seats_released: INTEGER - Seats allocated
- txn_status: TEXT - Payment status
- team_details: JSONB - Team registration data
- metadata: JSONB - Additional payment data
```

**Design Rationale**: 
- Supports both individual and team registrations
- JSONB fields accommodate flexible payment gateway responses
- GST tracking for financial compliance

#### `teams` & `team_members`
Team structure for group events.

```sql
Teams:
- id: UUID (PK)
- team_name: CITEXT - Case-insensitive unique names per event
- event_id: UUID (FK)
- leader_name: TEXT
- booking_id: UUID (FK)
- metadata: JSONB

Team Members:
- id: UUID (PK)
- team_id: UUID (FK)
- student_id: UUID (FK)
- student_role: TEXT - Role in team
- student_name/email: TEXT - Cached participant info
```

**Design Rationale**: 
- CITEXT ensures case-insensitive team name uniqueness
- Cached student info preserves data even if student records change

#### `solo_event_participant` & `team_events_attendance`
Attendance tracking for solo and team events with check-in/check-out timestamps.

### User Interaction Tables

#### `favourites`
Event bookmarking system for students.

#### `dispute`
Transaction dispute resolution system.

```sql
Columns:
- id: UUID (PK)
- txn_id: TEXT (FK) - Related booking
- student_email: TEXT (FK) - Dispute initiator
- description: TEXT - Dispute details
- event_id: UUID (FK) - Related event
- dispute_status: ENUM(OPEN, CLOSED_AS_TRUE, CLOSED_AS_FALSE)
```

### Accommodation System

#### `hostel_metadata`
Accommodation facility management.

```sql
Columns:
- id: UUID (PK)
- hostel_name: TEXT - Building name
- room_count: INTEGER - Total capacity
- is_male: BOOLEAN - Gender designation
- warden_email/password: TEXT - Staff credentials
- latitude/longitude: TEXT - GPS coordinates
- map_url: TEXT - Location link
- amrita_dayscholar_price: INTEGER - Internal rate
- non_amrita_price: INTEGER - External rate
- room_filled: INTEGER - Current occupancy
```

**Design Rationale**: Differential pricing for internal vs external participants.

#### `accomodation_personell`
Hospitality staff with independent authentication.

#### `accomodation_details`
Student accommodation bookings with payment tracking.

```sql
Columns:
- id: UUID (PK)
- student_id: UUID (FK)
- hostel_id: UUID (FK)
- personal_info: name, email, phone_number
- accommodation_details: is_male, is_hosteller, college info
- room_preference: TEXT - Room type preferences
- is_amrita_campus: BOOLEAN - Campus resident status
- payment_status: TEXT - Payment tracking
- payment_expires: TIMESTAMP - Payment deadline
- day_count: INTEGER - Duration of stay
- check_in/check_out: TIMESTAMP - Stay period
```

#### `hostel_check_in`
Actual check-in/out tracking with staff audit trail.

#### `gate_management`
Campus access control with direction tracking.

```sql
Columns:
- id: UUID (PK)
- personell_id: UUID (FK) - Staff who processed entry
- student_id: UUID (FK) - Student entering/exiting
- direction: ENUM(IN, OUT) - Movement direction
- logged_at: TIMESTAMP - Entry/exit time
```

**Design Rationale**: Complete audit trail for security and attendance tracking.

## Key Design Decisions

### 1. UUID Primary Keys
- Prevents enumeration attacks
- Enables distributed ID generation
- Better for microservice architecture

### 2. Separated Authentication and Business Logic
- Temporary tables (`student_onboarding`, `password_reset`) isolate auth flows
- Refresh tokens stored in main user tables for session management

### 3. Flexible Metadata with JSONB
- `team_details` and `metadata` in bookings accommodate varying payment gateway responses
- Enables schema evolution without migrations

### 4. Comprehensive Audit Trails
- `created_at`/`updated_at` on all major tables
- Staff tracking in hospitality and gate management
- Transaction status progression

### 5. Role-Based Data Separation
- Different personnel types with separate authentication
- Clear separation between student data and operational data

### 6. Event Flexibility
- Multi-day support through `event_schedule`
- Team and solo event handling
- Online/offline mode support

## Index Strategy

- Unique constraints on natural keys (email, roll numbers)
- Partial indexes for optional fields (roll numbers)
- Composite indexes for common query patterns (email + event_id)

## Data Integrity

- Foreign key constraints with RESTRICT/UPDATE CASCADE
- CHECK constraints for business logic (room_filled <= room_count)
- ENUM types for controlled vocabularies

This database design supports a comprehensive event management platform with robust user management, flexible event structures, comprehensive payment tracking, and full hospitality management capabilities.