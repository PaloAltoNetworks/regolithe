# Model
model:
  rest_name: user
  resource_name: users
  entity_name: User
  package: todo-list
  description: Represent a user.
  aliases:
  - usr
  get: true
  update: true
  delete: true
  extends:
  - '@base'

# Attributes
attributes:
- name: firstName
  description: The first name
  type: string
  exposed: true
  stored: true
  required: true
  filterable: true
  format: free
  orderable: true

- name: lastName
  description: The last name
  type: string
  exposed: true
  stored: true
  required: true
  filterable: true
  format: free
  orderable: true

- name: userName
  description: the login
  type: string
  exposed: true
  stored: true
  required: true
  filterable: true
  format: free
  orderable: true
