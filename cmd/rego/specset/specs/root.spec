# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: root
  description: root object.
  get: true
  root: true

# Relations
relations:
- rest_name: object
  descriptions:
    create: Creates a new object.
    get: Retrieves the list of objects.
  get: true
  create: true
