# Model
model:
  rest_name: user
  resource_name: users
  entity_name: User
  package: todo-list
  group: core
  description: Represent a user.
  aliases:
  - usr
  get:
    description: Retrieves the user with the given ID.
  update:
    description: Updates the user with the given ID.
  delete:
    description: Deletes the user with the given ID.
    parameters:
      required:
      - - - confirm
      entries:
      - name: confirm
        description: this is required.
        type: boolean
  extends:
  - '@base'
  validations:
  - $username

# Attributes
attributes:
  v1:
  - name: archived
    description: the object is archived and not deleted.
    type: boolean
    exposed: true
    stored: true
    example_value: false
    getter: true
    setter: true

  - name: firstName
    description: The first name.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: firstName
    filterable: true
    orderable: true

  - name: lastName
    description: The last name.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: lastName
    filterable: true
    orderable: true

  - name: userName
    description: the login.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: userName
    filterable: true
    orderable: true
