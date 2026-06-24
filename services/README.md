# Services

Various system-level services written in Rust that are used by the application. 

These currently include:

* Caching
    * Caches data from APIs to improve performance and reduce network usage.
    ![Rust caching service sequence diagram](https://dc8hq8aq7pr04.cloudfront.net/rust-caching-service.png)

* Aggregation
    * Fetches each PokéAPI resource once and returns one structured JSON profile for the Go CLI to render.
    ![Rust aggregation service sequence diagram](https://dc8hq8aq7pr04.cloudfront.net/rust-aggregation-service.png)
    