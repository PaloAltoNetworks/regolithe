# Model
model:
  rest_name: list
  resource_name: lists
  entity_name: List
  package: todo-list
  description: Represent a a list of task to do.
  aliases:
  - lst
  get:
    description: Retrieves the list with the given ID.
    parameters:
    - name: lgp1
      description: this is lgp1.
      type: string
    - name: lgp2
      description: this is lgp2.
      type: boolean
  update:
    description: Updates the list with the given ID.
    parameters:
    - name: lup1
      description: this is lpu1.
      type: string
    - name: lup2
      description: this is lpu2.
      type: boolean
  delete:
    description: Deletes the list with the given ID.
    parameters:
    - name: ldp1
      description: this is ldp1.
      type: string
    - name: ldp2
      description: this is ldp2.
      type: boolean
  extends:
  - '@base'

# Attributes
attributes:
  v1:
  - name: creationOnly
    description: This attribute is creation only.
    type: string
    exposed: true
    stored: true
    creation_only: true
    filterable: true
    format: free
    orderable: true

  - name: date
    description: The date.
    type: time
    exposed: true
    stored: true
    filterable: true
    orderable: true

  - name: description
    description: The description.
    type: string
    exposed: true
    stored: true
    filterable: true
    format: free
    orderable: true

  - name: name
    description: The name.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: the name
    filterable: true
    format: free
    getter: true
    setter: true
    orderable: true

  - name: readOnly
    description: This attribute is readonly.
    type: string
    exposed: true
    stored: true
    read_only: true
    filterable: true
    format: free
    orderable: true

  - name: slice
    description: this is a slice.
    type: list
    exposed: true
    subtype: string
    stored: true
    filterable: true
    orderable: true

  - name: unexposed
    description: This attribute is not exposed.
    type: string
    stored: true
    filterable: true
    format: free
    orderable: true

# Relations
relations:
- rest_name: task
  get:
    description: yeye.
    parameters:
    - name: ltgp1
      description: this is ltgp1.
      type: string
    - name: ltgp2
      description: this is ltgp2.
      type: boolean
  create:
    description: yoyo.
    parameters:
    - name: ltcp1
      description: this is ltcp1.
      type: string
    - name: ltcp2
      description: this is ltcp2.
      type: boolean

- rest_name: user
  get:
    description: yeye.
