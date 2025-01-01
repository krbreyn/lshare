lshare (for local share) is a quick and dirty CLI tool for transferring files over a local network using HTTP.

This is a work in progress.

Quick notes:
- The file server object should be passed in an FileStore and an EndpointStore(?) which are responsible for registering, deleting, loading, and saving endpoints and files.
- use either flags or cobra/viper libraries for the cli
- optional crude encryption through password-locked archives
- instead of having the endpoint be the filename, randomly generate an all-lowercase alphanumerical passcode
- be able to select a directory or a grouping of files to zip and send
- ability to save, load, and then list endpoints for repeated/scheduled serving of the same files
