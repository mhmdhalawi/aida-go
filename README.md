# Go Project

## Installation

**Step 1: Clone the repository**

```bash
git clone https://github.com/mhmdhalawi/aida-go.git
```

**Step 2: Go to the project folder**

```bash
cd aida-go
```

**Optional Step: Replace or add your own JSON files**

- The project reads all .json files from the data folder.

- You can add new JSON files or replace existing ones to test with your own data.

```bash
[
  {
    "firstName" : "John",
    "lastName" : "Doe",
    "birthday" : "1990-05-14T00:00:00Z",
    "address" : "123 Main Street, New York, NY",
    "phoneNumber" : "+1-555-1234"
  }
]
```

**Step 3: Run the development server**

```bash
go run main.go
```
