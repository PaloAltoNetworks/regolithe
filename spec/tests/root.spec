# Model
model:
  rest_name: root
  resource_name: root
  entity_name: Root
  package: todo-list
  group: core
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
      entries:
      - name: rlgmp1
        description: this is rlgmp1.
        type: string
        example_value: rlgmp1

      - name: rlgmp2
        description: this is rlgmp2.
        type: boolean
        example_value: "true"
  create:
    description: you.
    parameters:
      entries:
      - name: rlcp1
        description: this is rlcp1.
        type: string
        example_value: rlcp1

      - name: rlcp2
        description: this is rlcp2.
        type: boolean
        example_value: "true"

- rest_name: user
  get:
    description: yey.
    parameters:
      entries:
      - name: rugmp1
        description: this is rugmp1.
        type: string
        example_value: rugmp1

      - name: rugmp2
        description: this is rugmp2.
        type: boolean
        example_value: "true"
  create:
    description: you.
    parameters:
      entries:
      - name: rucp1
        description: this is rucp1.
        type: string
        example_value: rucp1

      - name: rucp2
        description: this is rucp2.
        type: boolean
        example_value: "true"
