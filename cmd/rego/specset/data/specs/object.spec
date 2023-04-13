# Model
model:
  rest_name: object
  resource_name: objects
  entity_name: Objects
  package: default
  group: core
  description: This is random object.
  get:
    description: Gets the object.
  update:
    description: Updates the object.
  delete:
    description: Deletes the object.
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
