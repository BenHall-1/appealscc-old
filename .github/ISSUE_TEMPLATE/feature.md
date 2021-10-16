name: Feature Request
description: File a Feature Request
title: "[FEATURE]: "
labels: ["feature request", "triage"]
assignees:
  - BenHall-1
body:
  - type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this feature request!
  - type: textarea
    id: what-happened
    attributes:
      label: What are you expecting to happen?
      description: Also tell us, how do you see this feature being implemented?
      placeholder: I want a button that...
      value: "A bug happened!"
    validations:
      required: true
  - type: dropdown
    id: apis
    attributes:
      label: What section of the system does this feature affect?
      multiple: true
      options:
        - API
        - Website
        - Both
  - type: checkboxes
    id: terms
    attributes:
      label: Code of Conduct
      description: By submitting this issue, you agree to follow our Code of Conduct
      options:
        - label: I agree to follow this project's Code of Conduct
          required: true
