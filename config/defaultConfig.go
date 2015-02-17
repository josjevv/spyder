package config

const defaultYaml = `
---
mongohost: 127.0.0.1:27017
mongodb: safetyapps
components:
  history: true
  notifications: true
associations:
  all: [history, notifications]
`
