
# Calculate SLA





## Problem
 - Input date and time to create a reference SLA.
 - SLA time is calculated based on working days from Monday to Friday.
 - The working day starts at 09:00 AM and ends at 18:00 PM.
 - There is a break time from 12:00 PM to 13:00 PM.
 - SLA reference times are taken from a master source.

You want to calculate SLA deadlines at 50%, 75%, and 100%.

## Asumption I made

 - 12:00 PM - 14:00 PM counted as 1 hour ( 12:00 PM - 13:00 PM is break time, and 13:00 PM - 14:00 PM counted as working hour again)
 - 18:00 PM - 10:00 AM counted as 1 hour ( Outside Working hour until 9:00AM)
 - All second is floored


## API Reference

#### Calculate SLA

```http
  GET /calculate-sla
```

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `sla` | `string` | **Required**. Duration ( Ex: 24,48,72, etc ) |
| `start_time` | `string` | **Required**. Time Stamp ( Ex : 2023-09-04T09:00:00Z (it means 4 September 2023 on 9 AM) ) |



## How To Run

Clone this repo and run:

```bash
  go mod tidy
  go run main.go
```
access on localhost:8080