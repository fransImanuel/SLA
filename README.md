
# Calculate SLA





## Problem
 - Input the date and time to create a reference SLA.
 - SLA time is calculated based on working days from Monday to Friday.
 - The working day starts at 09:00 AM and ends at 18:00 PM.
 - There is a break time from 12:00 PM to 13:00 PM.
 - SLA reference times are taken from a master source.

You want to calculate SLA deadlines at 50%, 75%, and 100%.

## Assumption I made

 - 12:00 PM - 14:00 PM is counted as 1 hour (12:00 PM - 13:00 PM is break time, and 13:00 PM - 14:00 PM is counted as working hour again).
 - 18:00 PM - 10:00 AM is counted as 1 hour (outside working hours until 9:00 AM).
 - All seconds are floored.


## API Reference

#### Calculate SLA

```HTTP
  GET /calculate-sla
```

| Query Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `sla` | `string` | **Required**. Duration (e.g., 24, 48, 72, etc.) |
| `start_time` | `string` | **Required**. Timestamp (e.g., 2023-09-04T09:00:00Z, which means 4 September 2023 at 9 AM) |



## How To Run

Clone this repo and run:

```bash
  go mod tidy
  go run main.go
```
access on localhost:8080

[Link to SLA API Usage](https://docs.google.com/document/d/1EcA5Bq1tWTAAquR8SLItb9lkzqM3v7ZV7q9IDzAMvVM/edit)
