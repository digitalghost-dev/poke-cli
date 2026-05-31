# Services

Various system-level services written in Rust that are used by the application. 

These currently include:

* Caching
    * Caches data from APIs to improve performance and reduce network usage.

* Aggregation
    * Fetches each PokéAPI resource once and returns one structured JSON profile for the Go CLI to render.
    ![aggregation service squence diagram](https://dc8hq8aq7pr04.cloudfront.net/rust_aggregation_service.png)
    