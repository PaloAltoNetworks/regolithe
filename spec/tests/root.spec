# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: todo-list
  description: Root object of the API.
  get:
    description: Retrieve the root object.
  root: true

# Relations
relations:
- rest_name: list
  get:
    description: yey.
  create:
    description: you.

- rest_name: user
  get:
    description: yey.
  create:
    description: you.
