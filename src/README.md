# TODOs

- Augment membership list to support the virtual ring structure.
- connecting FD with sdfs
  - Add the channel write in membership list add/remove functions
  - Add the channel read in the passive_actions of sdfs
  - Might need a lock to block active actions when doing passive actions.
  - Maybe print an ack once this passive is done to hint the user that operations are now available.
- implement grpc functions
  - implement server side functions in active_actions in sdfs module.
  - implement client side functions in sds module.
  - Need to pay attention to the re-replication process. e.g. introducer must block until ack when joining multiple nodes.
  - This involves introducer to serialize the join requests.
