attributes:
- creation_only: true
  description: This attribute is creation only
  exposed: true
  filterable: true
  format: free
  name: creationOnly
  orderable: true
  stored: true
  type: string
- description: The date
  exposed: true
  filterable: true
  name: date
  orderable: true
  stored: true
  type: time
- description: The description
  exposed: true
  filterable: true
  format: free
  name: description
  orderable: true
  stored: true
  type: string
- description: The name
  exposed: true
  filterable: true
  format: free
  getter: true
  name: name
  orderable: true
  required: true
  setter: true
  stored: true
  type: string
  unique: true
- description: This attribute is readonly
  exposed: true
  filterable: true
  format: free
  name: readOnly
  orderable: true
  read_only: true
  stored: true
  type: string
- description: this is a slice
  exposed: true
  filterable: true
  name: slice
  orderable: true
  stored: true
  subtype: string
  type: list
- description: This attribute is not exposed
  filterable: true
  format: free
  name: unexposed
  orderable: true
  stored: true
  type: string
children:
- create: true
  get: true
  relationship: child
  rest_name: task
- get: true
  relationship: child
  rest_name: user
model:
  aliases:
  - lst
  delete: true
  description: Represent a a list of task to do.
  entity_name: List
  extends:
  - '@base'
  get: true
  package: todo-list
  resource_name: lists
  rest_name: list
  update: true
