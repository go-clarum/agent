### Clarum executor agent

### Lifecycle
At this point we always start the agent, initiate mocks, run test actions and after all that we shut it down.
Maybe in the future we will keep it alive with the mocks running to improve performance.

Keeping it alive would require:
- to manage the agent lifecycle
- add config flag to keep the agent alive
- introduce concept of TestRun
- split config in agent startup configuration & TestRun configuration
