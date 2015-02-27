package config

const defaultYaml = `
---
mongohost: 127.0.0.1:27017
mongodb: safetyapps
components:
  logger: true
  history: true
  notifications: true
associations:
  all: [logger, history, notifications]
notifications:
  all: true
`
