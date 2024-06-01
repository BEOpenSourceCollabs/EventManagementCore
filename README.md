# Event Management System

user management- login / sign up / roles assigned to users  ( organisers , admins and normal users) 

- Event organisers: Should be able to create , edit , update , delete event. 
- Event registration: Normal users can register for an event and since we have limited capacity for events (offline events) this will be checked. 
- Event notification: Once the user registration is complete they should get notified about those events 
- Administration: Admins should be able to deal with users , moderate events etc , see analytics about the events such as popularity/ views etc .
- Event reviews: Users can leave reviews for events they attended etc.

## Setup

```bash
git clone https://github.com/BEOpenSourceCollabs/EventManagementCore
cd EventManagementCore
```

### Run

```bash
make run
```

### Run Tests

```bash
make test

# or 

make test.unit
```