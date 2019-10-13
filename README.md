# API Serve
The purpose of `API Serve` is to be a very minimal yet production ready API server. `API Serve` 
was designed with web applications in mind that are primarily handing JSON request/responses. `API Serve`
supports global and route specific middleware. It also uses `httprouter` under 
the hood for request routing. The server supports load shedding out of the box for production
applications.
