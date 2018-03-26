# Model
model:
  rest_name: list
  resource_name: lists
  entity_name: List
  package: todo-list
  description: Represent a a list of task to do.
  aliases:
  - lst
  get: true
  update: true
  delete: true
  extends:
  - '@base'

# Attributes
attributes:
- name: creationOnly
  description: This attribute is creation only
  type: string
  exposed: true
  stored: true
  creation_only: true
  filterable: true
  format: free
  orderable: true

- name: date
  description: The date
  type: time
  exposed: true
  stored: true
  filterable: true
  orderable: true

- name: description
  description: The description
  type: string
  exposed: true
  stored: true
  filterable: true
  format: free
  orderable: true

- name: name
  description: The name
  type: string
  exposed: true
  stored: true
  required: true
  filterable: true
  format: free
  getter: true
  setter: true
  orderable: true
  example_value: "the name"

- name: readOnly
  description: This attribute is readonly
  type: string
  exposed: true
  stored: true
  read_only: true
  filterable: true
  format: free
  orderable: true

- name: slice
  description: this is a slice
  type: list
  exposed: true
  subtype: string
  stored: true
  filterable: true
  orderable: true

- name: unexposed
  description: This attribute is not exposed
  type: string
  stored: true
  filterable: true
  format: free
  orderable: true

# Relations
relations:
- rest_name: task
  descriptions:
    get: yeye
    create: yoyo
  get: true
  create: true

- rest_name: user
  descriptions:
     get: yeye
  get: true
