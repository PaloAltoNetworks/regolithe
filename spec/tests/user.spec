{
    "attributes": [
        {
            "description": "The first name",
            "exposed": true,
            "filterable": true,
            "format": "free",
            "name": "firstName",
            "orderable": true,
            "required": true,
            "stored": true,
            "type": "string"
        },
        {
            "description": "The last name",
            "exposed": true,
            "filterable": true,
            "format": "free",
            "name": "lastName",
            "orderable": true,
            "required": true,
            "stored": true,
            "type": "string"
        },
        {
            "description": "the login",
            "exposed": true,
            "filterable": true,
            "format": "free",
            "name": "userName",
            "orderable": true,
            "required": true,
            "stored": true,
            "type": "string",
            "unique": true
        }
    ],
    "model": {
        "aliases": [
            "usr"
        ],
        "delete": true,
        "description": "Represent a user.",
        "entity_name": "User",
        "extends": [
            "@base"
        ],
        "get": true,
        "package": "todo-list",
        "resource_name": "users",
        "rest_name": "user",
        "update": true
    }
}
