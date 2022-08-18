# Regolithe

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/670b3ed05f0c4d81b2215bf57b500672)](https://www.codacy.com/gh/PaloAltoNetworks/regolithe/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=PaloAltoNetworks/regolithe&amp;utm_campaign=Badge_Grade) [![Codacy Badge](https://app.codacy.com/project/badge/Coverage/670b3ed05f0c4d81b2215bf57b500672)](https://www.codacy.com/gh/PaloAltoNetworks/regolithe/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=PaloAltoNetworks/regolithe&amp;utm_campaign=Badge_Coverage)

> This README is a work in progress

Regolithe is a library and generation tool you can use to author and generate
data models from yaml-based specification. It is agnostic in term of output
language but provides out of the box a JSON Schema generator as well as a
Markdown documentation generator. The library exposes a function that lets you
configure a Cobra command you can integrate to your own generation tool. The
reference implementation is Elemental, that generates Go model for the Bahamut
Stack.

## Specification Set

Regolithe is basically a specification on how to write specifications, called
Specification Set. A Specification Set is a set of standalone Specification files
that are bound together by defining relationship.

For instance, we can imagine having:

* A Task specification file (task.spec)
* A List specification file (list.spec)
* A relationship allowing to get Task from a List.

Here's an example specification for the Item:

```yaml
model:
  rest_name: task
  resource_name: tasks
  entity_name: Task
  package: example
  group: core
  description: Task you have to do.
  get:
    description: Retrieves the Task with the given ID.
  update:
    description: Updates the Task with the given ID.
  delete:
    description: Deletes the Task with the given ID.
  extends:
  - '@identifiable'

# Ordering
default_order:
- :no-inherit
- date

# Indexes
indexes:
- - status

# Attributes
attributes:
  v1:
  - name: description
    description: Description of the Task.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: Buy milk

  - name: status
    description: The status of the Task
    type: enum
    exposed: true
    stored: true
    allowed_choices:
    - TODO
    - DONE
    - HOLD
    default_value: TODO
```

And here's an example specification for a List holding them:

```yaml
model:
  rest_name: list
  resource_name: lists
  entity_name: List
  package: example
  group: core
  description: List of Tasks.
  get:
    description: Retrieves the List with the given ID.
  update:
    description: Updates the List with the given ID.
  delete:
    description: Deletes the List with the given ID.
  extends:
  - '@identifiable'

# Attributes
attributes:
  v1:
  - name: name
    description: Name of the list.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: Buy milk

# Relations
relations:
- rest_name: task
  get:
    description: Retrieve the tasks in the list
    parameters:
      entries:
      - name: status
        description: Only get tasks with a specific status.
        type: enum
        allowed_choices:
        - TODO
        - DONE
        - HOLD

  create:
    description: Create a new task in the list.
```

## Specification sections

## Model

> TODO: describe the model section of a spec file
> TODO: describe detached specs

## Attributes

> TODO: describe the attributes properties

## Relationship

> TODO: describe how to create hierarchy

## Type Mappings

> TODO: describe how to map external types to an attribute

## Validation Mappings

> TODO: describe how to create attribute and full model validation functions
