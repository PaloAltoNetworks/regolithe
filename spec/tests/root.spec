children:
- create: true
  get: true
  relationship: root
  rest_name: list
- create: true
  get: true
  relationship: root
  rest_name: user
model:
  description: Root object of the API
  entity_name: Root
  get: true
  package: todo-list
  resource_name: root
  rest_name: root
  root: true
