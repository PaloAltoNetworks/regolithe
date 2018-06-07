# Model
model:
  rest_name: object
  resource_name: objects
  entity_name: Objects
  package: default
  description: This is random object.
  get: true
  update: true
  delete: true
  extends:
  - '@identifiable'

# Attributes
attributes:
  v1:
  - name: name
    description: The name of the object.
    type: string
    exposed: true
    stored: true
    filterable: true
    orderable: true
