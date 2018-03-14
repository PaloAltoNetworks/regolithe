{
    "attributes": [
        {
            "description": "The description",
            "exposed": true,
            "filterable": true,
            "format": "free",
            "name": "description",
            "orderable": true,
            "stored": true,
            "type": "string"
        },
        {
            "description": "The name",
            "exposed": true,
            "filterable": true,
            "format": "free",
            "getter": true,
            "name": "name",
            "orderable": true,
            "required": true,
            "setter": true,
            "stored": true,
            "type": "string"
        },
        {
            "allowed_choices": [
                "DONE",
                "PROGRESS",
                "TODO"
            ],
            "default_value": "TODO",
            "description": "The status of the task",
            "exposed": true,
            "filterable": true,
            "name": "status",
            "orderable": true,
            "stored": true,
            "type": "enum"
        }
    ],
    "model": {
        "aliases": [
            "tsk"
        ],
        "delete": true,
        "description": "Represent a task to do in a listd",
        "entity_name": "Task",
        "extends": [
            "@base"
        ],
        "get": true,
        "package": "todo-list",
        "resource_name": "tasks",
        "rest_name": "task",
        "update": true
    }
}
