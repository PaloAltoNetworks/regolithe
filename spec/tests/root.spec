# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: todo-list
  description: Root object of the API.
  get:
    description: Retrieve the root object.
  root: true

# Relations
relations:
- rest_name: list
  get:
    description: yey.
    parameters:
    - name: rlgmp1
      description: this is rlgmp1.
      type: string
    - name: rlgmp2
      description: this is rlgmp2.
      type: boolean
  create:
    description: you.
    parameters:
    - name: rlcp1
      description: this is rlcp1.
      type: string
    - name: rlcp2
      description: this is rlcp2.
      type: boolean

- rest_name: user
  get:
    description: yey.
    parameters:
    - name: rugmp1
      description: this is rugmp1.
      type: string
    - name: rugmp2
      description: this is rugmp2.
      type: boolean
  create:
    description: you.
    parameters:
    - name: rucp1
      description: this is rucp1.
      type: string
    - name: rucp2
      description: this is rucp2.
      type: boolean
