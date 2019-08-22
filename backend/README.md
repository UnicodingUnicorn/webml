# backend

backend for WebML. Mostly acts as a proxy for [minio](min.io).

## API

| Model           | Parser           | Model Data (Data)    | Model Data (Labels)    |
| --------------- | ---------------- | -------------------- | ---------------------- |
| Get Models      | Get Parsers      | Get Model Data       | Get Model Labels       |
| Get Model by ID | Get Parser by ID | Get Model Data By ID | Get Model Labels By ID |
| Upload Model    | Upload Parser    | Upload Model Data    | Upload Model Labels    |

| Batch            |  Session        |
| ---------------- | --------------- |
| Get Batch        | Get Loss        |
| Get Batch Data   | Post Loss       |
| Get Batch Labels | Update Weights  |
| Batch Data       | New Session     |

Most of the routes, save for the ones under `Session`, are basically just proxy routes that redirect one to a minio pre-signed URL.

### Model Handlers

#### Get Models

```
GET /models
```

Get a list of all the model IDs available.

##### Success (200 OK)

```json
[ "<model 1 id>", "<model 2 id>", "..." ]
```

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

#### Get Model By ID

```
GET /model/:id
```

Get a specific model's definition by ID.

##### URL Params

| Name | Type   | Description   |
| ---- | ------ | ------------- |
| id   | String | Model's ID.   |

##### Success

Redirect to minio pre-signed URL to download the data.

##### Errors

| Code | Description                                             |
| ---- | ------------------------------------------------------- |
| 404  | The model referenced by the supplied id does not exist. |
| 500  | Error querying minio.                                   |

---

#### Upload Model

```
PUT /model
```

Upload a model.

##### Body (multipart/form-data)

| Name | Type   | Description          | Required |
| ---- | ------ | -------------------- | -------- |
| id   | String | ID of the new model. | ✓        |

##### Success

Redirect to minio pre-signed URL to upload the data.

##### Errors

| Code | Description                 |
| ---- | --------------------------- |
| 400  | Error parsing request body. |
| 500  | Error querying minio.       |

---

### Parser Definitions

Parsers are lua scripts, written by the user to parse specific forms of data into tensors.

#### Get Parsers

```
GET /parsers
```

Get a list of all the parser IDs available.

##### Success (200 OK)

```json
[ "<parser 1 id>", "<parser 2 id>", "..." ]
```

---

#### Get Parser By ID

```
GET /parser/:id
```

Get a specific parser's definition by ID.

##### URL Params

| Name | Type   | Description   |
| ---- | ------ | ------------- |
| id   | String | Parser's ID.  |

##### Success

Redirect to minio pre-signed URL to download the data.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

#### Upload Parser

```
PUT /parser
```

Upload a parser.

##### Success

Redirect to minio pre-signed URL to upload the data.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

### Model Data Handlers

These routes handle the training data associated with the model. Training data is split into two - the data and the labels. Each set of data must have an associated set of labels. Both need to share the same ID.

#### Get Model Data

```
GET /model/:model/data
```

Get a list of all the data IDs available for a model.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |

##### Success (200 OK)

```json
[ "<model data 1 id>", "<model data 2 id>", "..." ]
```

---

#### Get Model Data By ID

```
GET /model/:model/data/:id
```

Get a specific model data's definition by ID.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |
| id    | String | Model Data's ID.          |

##### Success

Redirect to minio pre-signed URL to download the model data.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

#### Upload Model Data

```
PUT /model/:model/data
```

Upload data to a model.

##### URL Params

| Name  | Type   | Description                   |
| ----- | ------ | ----------------------------- |
| model | String | ID of the model to upload to. |

##### Body (multipart/form-data)

| Name | Type   | Description            | Required |
| ---- | ------ | ---------------------- | -------- |
| id   | String | ID of the new dataset. | ✓        |

##### Success

Redirect to minio pre-signed URL to upload the data.

##### Errors

| Code | Description                 |
| ---- | --------------------------- |
| 400  | Error parsing request body. |
| 500  | Error querying minio.       |

---

#### Get Model Labels

```
GET /model/:model/labels
```

Get a list of all the label IDs available for a model.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |

##### Success (200 OK)

```json
[ "<model labels 1 id>", "<model labels 2 id>", "..." ]
```

---

#### Get Model Data By ID

```
GET /model/:model/labels/:id
```

Get a specific model label's definition by ID.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |
| id    | String | Model Label's ID.         |

##### Success

Redirect to minio pre-signed URL to download the model labels.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

#### Upload Model Data

```
PUT /model/:model/labels
```

Upload labels to a model.

##### URL Params

| Name  | Type   | Description                   |
| ----- | ------ | ----------------------------- |
| model | String | ID of the model to upload to. |

##### Body (multipart/form-data)

| Name | Type   | Description             | Required |
| ---- | ------ | ----------------------- | -------- |
| id   | String | ID of the new labelset. | ✓        |

##### Success

Redirect to minio pre-signed URL to upload the labels.

##### Errors

| Code | Description                 |
| ---- | --------------------------- |
| 400  | Error parsing request body. |
| 500  | Error querying minio.       |

---

### Batch Handlers

These routes handle the separation of model data/labelsets into small "batches" for training by individual nodes.

#### Get Batch

```
GET /model/:model/batch
```

Get the ID of a random batch from a model.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |

##### Success (200 OK)

String ID of the batch.

---

#### Get Batch Data by ID

```
GET /model/:model/batch/:id/data
```

Get a batched dataset of a model by ID.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |
| id    | String | ID of the batch to query. |

##### Success

Redirect to minio pre-signed URL to download the batch of model data.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

#### Get Batch Labels by ID

```
GET /model/:model/batch/:id/labels
```

Get a batched labelset of a model by ID.

##### URL Params

| Name  | Type   | Description               |
| ----- | ------ | ------------------------- |
| model | String | ID of the model to query. |
| id    | String | ID of the batch to query. |

##### Success

Redirect to minio pre-signed URL to download the batch of model labels.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

#### Batch Data

```
POST /model/:model/data/:id/batch
```

Run batching algorithm on the data of a model.

##### URL Params

| Name  | Type   | Description                 |
| ----- | ------ | --------------------------- |
| model | String | ID of the model to query.   |
| id    | String | ID of the dataset to batch. |

##### Body (multipart/form-data)

| Name         | Type   | Description                        | Required |
| ------------ | ------ | ---------------------------------- | -------- |
| data_parser  | String | ID of the parser for the dataset.  | ✓        |
| label_parser | String | ID of the parser for the labelset. | ✓        |
| batch_size   | int    | Number of data-things in a batch.  | ✓        |

##### Success

Empty body.

##### Errors

| Code | Description           |
| ---- | --------------------- |
| 500  | Error querying minio. |

---

### Session Handlers

These routes handle the sessions in which a model is trained. Usually, one only has one session per model.

#### Get Loss

```
GET /session/:id/loss
```

Get the loss associated with a specific session.

##### URL Params

| Name | Type   | Description   |
| ---- | ------ | ------------- |
| id   | String | Session's ID. |

##### Success (200 OK)

Value of the loss of the session (float).

##### Errors

| Code | Description                                      |
| ---- | ------------------------------------------------ |
| 404  | Session with the supplied ID could not be found. |

---

#### Post Loss

```
POST /session/:id/loss
```

Update the loss associated with a specific session.

##### URL Params

| Name | Type   | Description   |
| ---- | ------ | ------------- |
| id   | String | Session's ID. |

##### Body (application/json)

| Name | Type  | Description     | Required |
| ---- | ----- | --------------- | -------- |
| loss | float | New loss value. | ✓        |

##### Success (200 OK)

Empty body.

##### Errors

| Code | Description                                      |
| ---- | ------------------------------------------------ |
| 400  | Error parsing the input body.                    |
| 404  | Session with the supplied ID could not be found. |

---

#### Update Weights

```
POST /session/:id/weights
```

Update the weights of a session. This performs a weighted average of one's new values and the existing ones based on the `alpha` set when the session was created.

##### URL Params

| Name | Type   | Description   |
| ---- | ------ | ------------- |
| id   | String | Session's ID. |

##### Body (application/json)

| Name  | Type    | Description                  | Required |
| ----- | ------- | ---------------------------- | -------- |
| shape | []int   | Shape of the weights tensor. | ✓       |
| data  | []float | Data of the weights tensor.  | ✓       |

##### Success (200 OK)

Empty body.

##### Errors

| Code | Description                                                                          |
| ---- | ------------------------------------------------------------------------------------ |
| 400  | Error parsing the input body/Tensor shape or length does not match the existing one. |
| 404  | Session with the supplied ID could not be found.                                     |

---

#### New Session

```
POST /session/:id
```

Create a new session. The onus is on the client to at least make a decent effort of making sure the id is unique.

##### URL Params

| Name | Type   | Description                          |
| ---- | ------ | ------------------------------------ |
| id   | String | The ID of the session to be created. |

##### Body (application/json)

| Name  | Type  | Description                   | Required | Default |
| ----- | ----- | ----------------------------- | -------- | ------- |
| shape | []int | Shape of the weights tensor.  | ✓       | -        |
| loss  | float | Initial loss of the network.  | X       | 0.0      |
| alpha | float | Alpha of the network. New weights are integrated according to the formula `old * ALPHA + new * (1 - ALPHA)`. |  X       | 0.9      |

##### Success (200 OK)

Empty body.

##### Errors

| Code | Description                                  |
| ---- | -------------------------------------------- |
| 400  | Error parsing the input body.                |
| 409  | Session with the supplied ID already exists. |