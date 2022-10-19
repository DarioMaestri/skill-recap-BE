# Progetto per annotare skill acquisite

## Scopo personale

Lo scopo del progetto era imparare le basi del linguaggio di programmazione GoLang e ripassare il FE Angular.

## Idea del programma

Lo scopo del programma è quello di salvare le "skill/compentenze" dell'utente.

## Stato FE e BE

Il FE non è finito, esiste solo una bozza di login e di come dovrebbe apparire, ma il mio scopo era ripassare Angular e flexbox (motivo per cui ci sono i bordi marcati di alcune sezioni).
Il BE invece lo considero finito, con la possibilità di aggiungere, rimuovere e modificare Skill e Utenti tramite chiamate REST. Attualmente si connette ad un db mySQL in locale.

## Run

go run main.go

## Chiamate di esempio per testare

I Body di request e response sono in formato JSON.
In caso di errore la response sarà del tipo:
404 Not Found / 500 KO

```json
{
	"error": "Motivo dell'errore"
}
```

### Skill /skills

#### GetSkills

Recupera tutte le skill a DB.
GET /skills
Response:

```json
[
	{
		"id": 4,
		"name": "Angular",
		"version": "6"
	},
	{
		"id": 3,
		"name": "Angular",
		"version": "JS"
	},
	{
		"id": 5,
		"name": "IntelliJ IDEA",
		"version": ""
	},
	{
		"id": 2,
		"name": "Java",
		"version": "1.11"
	},
	{
		"id": 1,
		"name": "Java",
		"version": "1.8"
	}
]
```

#### GetSkill

Recupera solo la skill con {id} passato nell'url
GET /skills/{id}
GET /skills/3
Response:

```json
{
	"id": 3,
	"name": "Angular",
	"version": "JS"
}
```

#### InsertSkill

Aggiunge una nuova skill
POST /skills
Body:

```json
{
	"name": "ABC",
	"version": "4.6"
}
```

Response:
200 OK - {id} - Id con cui e' stato salvato

#### DeleteSkill

Elimina la skill con {id} passato nell'url
DELETE /skills/{id}
DELETE /skills/6

#### UpdateSkill

Aggiorna la skill con {id} passato nell'url
PUT /skills/{id}
PUT /skills/6
Body:

```json
{
	"name": "ABC",
	"version": "4.4"
}
```

### User /users

#### Login

Effettua la login
POST /users/login
Body:

```json
{
	"username": "admin",
	"password": "password"
}
```

Response: boolean

#### GetUsers

Recupera tutti gli utenti a DB.
GET /users
Response:

```json
[
	{
		"id": 1,
		"username": "admin",
		"password": "password"
	},
	{
		"id": 2,
		"username": "dario",
		"password": "password"
	},
	{
		"id": 3,
		"username": "dario2",
		"password": "dario"
	}
]
```

#### GetUser

Recupera solo l'utente con {id} passato nell'url
GET /users/{id}
GET /users/3
Response:

```json
{
	"id": 3,
	"username": "dario2",
	"password": "dario"
}
```

#### InsertUser

Aggiunge un nuovo user
POST /users
Body:

```json
{
	"username": "ABC",
	"password": "4.6"
}
```

Response:
200 OK - {id} - Id con cui e' stato salvato

#### DeleteUser

Elimina l'utente con {id} passato nell'url
DELETE /users/{id}
DELETE /users/6

#### UpdateUser

Aggiorna l'utente con {id} passato nell'url
PUT /users/{id}
PUT /users/6
Body:

```json
{
	"username": "ABC",
	"password": "4.4"
}
```

### UserSkill /userskill

#### GetUserSkills

Recupera tutte le associazioni utente skill
GET /userskill
Response:

```json
[
	{
		"userId": 2,
		"skillId": 1
	},
	{
		"userId": 1,
		"skillId": 3
	},
	{
		"userId": 2,
		"skillId": 4
	}
]
```

#### GetUsersBySkill

Recupera tutti gli utenti che hanno {skillId}
GET /userskill/users/{skillId}
GET /userskill/users/1
Response:

```json
[
	{
		"id": 2,
		"username": "dario",
		"password": ""
	}
]
```

#### GetSkillsByUsers

Recupera tutte le skill per {userId}
GET /userskill/skills/{userId}
GET /userskill/skills/2
Response:

```json
[
	{
		"id": 2,
		"name": "Java",
		"version": "1.11"
	},
	{
		"id": 2,
		"name": "Java",
		"version": "1.11"
	}
]
```

#### InserUserSkill

Aggiunge una skill all'utente
POST /userskill
Body :

```json
{
	"user_id": 2,
	"skill_id": 3
}
```

#### DeleteUserSkill

Elimina l'associazione della {skillId} all' {userId}
DELETE /userskill/{userId}/{skillId}
