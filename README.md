# Polls

Polls is a public JSON REST API that lets you create, view, edit, delete and vote on polls. Voting is limited to one vote per IP address per poll.
Polls can be configured by setting expiry time, results visibility and private status. The API also provides the ability to list, search and sort public polls.

## How to run locally

Requirements:

- Bash
- Docker

1. `git clone https://github.com/ivcp/polls.git`
2. `cd polls`
3. create a `.env` file in the repository's root directory (see `.env.example`)
4. make sure Docker is running
5. `bash build.sh`
6. `curl localhost/v1/healthcheck` to check if it's working

## API Usage

### POST /v1/polls

Creates new poll. It's necessary to provide a question and at least two options. Option positions must also be provided and start at 0.

Example request body:

```
{
  "question": "Favourite color?",
  "options": [
    { "value": "Red", "position": 0 },
    { "value": "Blue", "position": 1 }
  ]
}
```

Optionally you can provide:

- `"description"` - poll description.
- `"expires_at"` - time when the poll expires. Must be at least two minutes in the future. [ISO 8601](https://www.iso.org/iso-8601-date-and-time-format.html) string e.g. "2024-02-05T14:48:00.000Z".
- `"is_private"` - private polls are only accessible by link.
- `"results_visibility"` - when results can be seen. Accepted values: "always", "after_vote", "after_deadline".

<details>
  <summary>Example response:</summary>

```
{
"poll": {
  "id": "6df661aa-4f3f-4281-8b69-da430a8ebad4",
  "question": "Favourite color?",
  "description": "",
  "options": [
    {
      "id": "802c593f-5f79-44f7-80d1-4cc4e40ddcec",
      "value": "Red",
      "position": 0
    },
    {
      "id": "8ea93888-8002-4889-94a1-24d75e10c07d",
      "value": "Blue",
      "position": 1
    }
  ],
  "created_at": "2024-02-26T17:19:44Z",
  "updated_at": "2024-02-26T17:19:44Z",
  "expires_at": "",
  "results_visibility": "always",
  "is_private": false,
  "token": "ZLCQIKYQ4MT7K2NJCRQWC4KMMU"
}
}
```

</details>

### GET /v1/polls/{poll ID}

Show individual poll.

<details>
  <summary>Example response:</summary>

```
{
"poll": {
  "id": "6df661aa-4f3f-4281-8b69-da430a8ebad4",
  "question": "Favourite color?",
  "description": "",
  "options": [
    {
      "id": "802c593f-5f79-44f7-80d1-4cc4e40ddcec",
      "value": "Red",
      "position": 0
    },
    {
      "id": "8ea93888-8002-4889-94a1-24d75e10c07d",
      "value": "Blue",
      "position": 1
    }
  ],
  "created_at": "2024-02-26T17:19:44Z",
  "updated_at": "2024-02-26T17:19:44Z",
  "expires_at": "",
  "results_visibility": "always",
  "is_private": false
}
}
```

</details>

### GET /v1/polls

List public polls.

Accepts query parameters:

- `search` - search by question
- `page_size` - set number of results per page _(default 20)_
- `page` - set current page number _(default 1)_
- `sort` - sort by:
  - `-created_at` latest created _(default)_
  - `created_at` oldest created
  - `-question` poll question in descending alphabetical order
  - `question` poll question in ascending alphabetical order

<details>
  <summary>Example response:</summary>

```
  {
  "metadata": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 1,
    "total_records": 1
  },
  "polls": [
    {
      "id": "6df661aa-4f3f-4281-8b69-da430a8ebad4",
      "question": "Favourite color?",
      "description": "",
      "options": [
        {
          "id": "802c593f-5f79-44f7-80d1-4cc4e40ddcec",
          "value": "Red",
          "position": 0
        },
        {
          "id": "8ea93888-8002-4889-94a1-24d75e10c07d",
          "value": "Blue",
          "position": 1
        }
      ],
      "created_at": "2024-02-26T17:19:44Z",
      "updated_at": "2024-02-26T17:19:44Z",
      "expires_at": "",
      "results_visibility": "always",
      "is_private": false
    }
  ]
}
```

</details>

### POST /v1/polls/{poll ID}/options/{option ID}

Vote for option.

<details>
  <summary>Example response:</summary>

```
{
  "message":"vote successful"
}
```

</details>

### GET /v1/polls/{pollID}/results

Show results for poll.

<details>
  <summary>Example response:</summary>

```
{
  "results": [
    {
      "id": "802c593f-5f79-44f7-80d1-4cc4e40ddcec",
      "value": "Red",
      "position": 0,
      "vote_count": 0
    },
    {
      "id": "117d4ef6-322e-436c-9c6b-46964e10b8c3",
      "value": "Green",
      "position": 2,
      "vote_count": 0
    },
    {
      "id": "8ea93888-8002-4889-94a1-24d75e10c07d",
      "value": "Blue",
      "position": 1,
      "vote_count": 1
    }
  ]
}
```

</details>

<hr>

**Token is required for following endpoints.** Token is generated when a poll is created and must be included in the Authorization header.

### PATCH /v1/polls/{poll ID}

Update poll question, description or expiration time. Supports partial updates.

Example request body:

```
{
  "description": "We all know there are only two colors."
}
```

Headers example:

`Authorization: Bearer ZLCQIKYQ4MT7K2NJCRQWC4KMMU`

<details>
  <summary>Example response:</summary>

```
  {
  "poll": {
    "id": "6df661aa-4f3f-4281-8b69-da430a8ebad4",
    "question": "Favourite color?",
    "description": "We all know there are only two colors.",
    "options": [
      {
        "id": "802c593f-5f79-44f7-80d1-4cc4e40ddcec",
        "value": "Red",
        "position": 0
      },
      {
        "id": "8ea93888-8002-4889-94a1-24d75e10c07d",
        "value": "Blue",
        "position": 1
      }
    ],
    "created_at": "2024-02-26T17:19:44Z",
    "updated_at": "2024-02-26T19:11:00Z",
    "expires_at": "",
    "results_visibility": "always",
    "is_private": false
  }
}
```

</details>

### DELETE /v1/polls/{poll ID}

Delete a poll.

Headers example:

`Authorization: Bearer ZLCQIKYQ4MT7K2NJCRQWC4KMMU`

<details>
  <summary>Example response:</summary>

```
{
  "message": "poll successfully deleted"
}
```

</details>

### POST /v1/polls/{poll ID}/options

Add option to poll.

Example request body:

```
{
  "value":"Green"
}
```

Headers example:
`Authorization: Bearer ZLCQIKYQ4MT7K2NJCRQWC4KMMU`

<details>
  <summary>Example response:</summary>

```
{
  "message":"option added successfully"
}
```

</details>

### PATCH /v1/polls/{pollID}/options/{optionID}

Update an existing option's value.

Example request body:

```
{"value":"Yellow"}
```

Headers example:
`Authorization: Bearer ZLCQIKYQ4MT7K2NJCRQWC4KMMU`

<details>
  <summary>Example response:</summary>

```
{
  "message":"option updated successfully"
}
```

</details>

### PATCH /v1/polls/{pollID}/options

Update option positions.

Example request body:

```
{
  "options": [
    {
      "id": "802c593f-5f79-44f7-80d1-4cc4e40ddcec",
      "position": 2
    },
    {
      "id": "117d4ef6-322e-436c-9c6b-46964e10b8c3",
      "position": 0
    }
  ]
}
```

Headers example:
`Authorization: Bearer ZLCQIKYQ4MT7K2NJCRQWC4KMMU`

<details>
  <summary>Example response:</summary>

```
{
  "message":"option updated successfully"
}
```

</details>

### DELETE /v1/polls/{pollID}/options/{optionID}

Delete option.

Headers example:
`Authorization: Bearer ZLCQIKYQ4MT7K2NJCRQWC4KMMU`

<details>
  <summary>Example response:</summary>

```
{
  "message":"option deleted successfully"
}
```

</details>

## Technologies used:

- Go
- Postgres
- Docker
- Caddy

## Contributing

If you'd like to contribute, please fork the repository and open a pull request to the `master` branch.
