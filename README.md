lshare (for local share) is a quick and dirty CLI tool for transferring files over a local network using HTTP.

This is a work in progress.

Quick notes:
- The file server object should be passed in an FileStore and an EndpointStore(?) which are responsible for registering, deleting, loading, and saving endpoints and files.
