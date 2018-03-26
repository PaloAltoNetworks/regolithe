# Model
model:
  rest_name: task
  resource_name: tasks
  entity_name: Task
  package: todo-list
  description: Represent a task to do in a listd.
  aliases:
  - tsk
  get: true
  update: true
  delete: true
  extends:
  - '@base'

# Attributes
attributes:
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

- name: status
  description: The status of the task.
  type: enum
  exposed: true
  stored: true
  allowed_choices:
  - DONE
  - PROGRESS
  - TODO
  default_value: TODO
  filterable: true
  orderable: true
