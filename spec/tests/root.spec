# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: todo-list
  description: Root object of the API
  get: true
  root: true

# Relations
relations:
- rest_name: list
  get: true
  create: true

- rest_name: user
  get: true
  create: true
