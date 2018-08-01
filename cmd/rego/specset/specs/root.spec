# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: root
  description: root object.
  get:
    description: gets the object.
  root: true

# Relations
relations:
- rest_name: object
  get:
    description: Retrieves the list of objects.
  create:
    description: Creates a new object.
