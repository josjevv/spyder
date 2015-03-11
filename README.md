# spyder
Let's snoop all those data changes, shall we.

## Prerequisites
###go packages an other stuff
* gtm (see http://go-search.org/view?id=github.com%2Frwynn%2Fgtm)
* bazaar (bzr) needed to install gtm
* yaml (see https://github.com/go-yaml/yaml)

```shell
go get gopkg.in/mgo.v2
brew install bzr
go get github.com/rwynn/gtm
go get gopkg.in/yaml.v2
```

###enable replicaset in mongo
* close running mongo instance if needed
* restart mongo using right db paths etc using replSet

```shell
mongod --port 27017 --dbpath /data/db --replSet rs0
```

Connect to mongo
```shell
mongo
```

Initiate the replicaset and check for status
```mongo
rs.initiate()
rs.status()
```

### Dev path

* Pull the latest source code for API, Router & CarpetJs.

* Create a Mongo Oplog reader in golang. You can use http://go-search.org/view?id=github.com%2Frwynn%2Fgtm for reference.

* Clone the Spyder repository.

* Spyder should consume command line --yaml file. You can use https://github.com/go-yaml/yaml to parse the file. The yaml file should contain the following sections:

```yaml
---
components:
 [component_type]: [boolean]

associations:
 <collection_name>: <component>
 <collection_name2>: [<component1>, <component2>]
```

Example implementation is:

```yaml
---
components:
 history: true
 notifications: false
associations:
        incidents: [history, notifications]
```

The YAML configuration will decide what all collections we need to parse and what all components are associated with it. This helps us build an extensible and pluggable system. where turning off listeners is easy. Add a default for *, if present means all collections will abide by those unless specifically overridden.

On system upstart spyder must loop through all associations and generate an array of namespaced collections that need to be monitored for changes.

* Every Oplog entry should be matched, if the namespace is in the whitelist computed in the previous step. If not, Spyder should move on.

* If the entry is worthy of a change, add it to the listener's unbufferred channel. We will use Procuder/Sink/Dispatcher pattern to invoke the desired Listener.

* Create separate subpackages for each plugin. Each plugin must implement the Listener interface which has two bound struct methods

```
type Listener interface {
    new_event()
    submit_event()
}
```

* Listener Package can have its own structure underneath. Does not matter.

* Spyder should be resumable, if it resumes after a while it should know where it last left at. Maybe Spyder should have its own storage? https://github.com/HouzuoGuo/tiedot
* 

### API

1. Adding a new channel to notification setting

   **Method**: `POST`
   
  **EndPoint**: `/topic/<notification_ident>/channel/<channel_ident>`
  
  **Request body** :
  ```js
 {
      org:"",
      app_name: "",
      user:"",
      ident:""
  }
  ``` 
    **Response Code**: `200`
  
2. Removing a channel from notification setting

   **Method**: `DELETE`

  **EndPoint**: `/topic/<notification_ident>/channel/<channel_ident>`
  
  **Request body** :
  
  ```js
 {
      org:"",
      app_name: "",
      user:"",
      ident:""
  }
  ``` 
  **Response Code**: `200`

3. Removing  notification setting

  **Method** : `DELETE`
 
  **EndPoint**: `/topic/<notification_ident>`
 
  **Request Body**:
 
```js
 {
      org:"",
      app_name: "",
      user:"",
      ident:""
 }
```
  **Response Code**: `200`

4. Get all notification settings

 **Method** : `GET`
 
 **EndPoint**: `/topics`
 
 **Request Filteting**: `app_name` `org` `user`
 
 **Response Body**:
 
```js
 [
        {
            "_id": "",
            "created_on": 1425547531188,
            "updated_on": 1425879125700,
            "user": "",
            "org": "",
            "app_name": "",
            "channels": [
                "",
                ""
            ],
            "ident": ""
        }
]
```
**Response Code**: `200`

5. Get all channels

 **Method** : `GET`
 
 **EndPoint**: `/channels`
 
 **Request Filteting**: `app_name` `org` `user`
 
 **Response Body**:
 
```js
 [
        {
            "_id": "",
            "created_on": 1425545240236,
            "updated_on": 1425545240236,
            "user": "",
            "org": "",
            "app_name": "",
            "data": {
            },
            "ident": ""
        }
]
```
**Response Code**: `200`


**Note**
For all of the above request you must pass atleast one of the `org`, `app_name` and `user`.
`ident` denotes the notification type (e.g incident_title_changed)

