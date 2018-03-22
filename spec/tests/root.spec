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
  descriptions:
    get: yey
    create: you
  get: true
  create: true

- rest_name: user
  descriptions:
    get: yey
    create: you
  get: true
  create: true

