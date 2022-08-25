# Storage and update service
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/rog-golang-buddies/api-hub_storage-and-update-service/main.svg)](https://results.pre-commit.ci/latest/github/rog-golang-buddies/api-hub_storage-and-update-service/main)

## Description
This service represents API Specification document (ASD) data access interface (RabbitMQ to save/update, gRPC to get).
It provides gRPC API to search and get ASD documents, queue event listener to save ASD documents.
Also, this service is supposed to update ASD models. (request update from data scraping service)

### Main functions (To Do)
#### Proof of Concept stage
1. Provide gRPC API with methods:
    * search(search string, page int) (Page[ApiSpecDocShort], error) to search by domain and description. Return a list of short descriptions.
    * get(id long) (ApiSpecDef, error) to get the full API Spec document by id.
2. Process events from Rabbit MQ with ASD model:
    * retrieve events from Rabbit MQ with ASD model
    * save/update data in DB (set created_at/updated_at)
    * notify API Gateway on update if required (event need to contain flag - is notification needed)
3. Check last updated dates and request update if necessary:
    * get expiration time from the environment variable
    * perform periodically task and get records that require update
    * send them to data scraper service with queue

#### Future plans
1. Add caching.

### Data flows involving service
1. Usual UI request (search/get)
   API Gateway -gRPC-> storage and update service -> DB
2. Add new/update existing ASD
   data scraper service -RabbitMQ-> storage and update service (-RabbitMQ->, -RabbitMQ->) (API Gateway, DB)
3. Update stale data
   storage and update service -RabbitMQ-> data scraper service -RabbitMQ-> storage and update service -> DB

### Additional resources
1. Database schema
    https://dbdiagram.io/d/630287b1f1a9b01b0fb26570 / https://imgur.com/a/evcpXry
