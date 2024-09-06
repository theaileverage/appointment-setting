Services
Create services by creating Go packages and defining API endpoints within them. Read the docs to learn more.
Creating services
With Encore you create a service by defining an API within a regular Go package. Encore recognizes this as a service, and uses the package name as the service name.

On disk it might look like this:

/my-app
├── encore.app          // ... and other top-level project files
│
├── hello               // hello service (a Go package)
│   ├── hello.go        // hello service code
│   └── hello_test.go   // tests for hello service
│
└── world               // world service (a Go package)
    └── world.go        // world service code

Service Initialization
Under the hood Encore automatically generates a main function that initializes all your infrastructure resources when the application starts up. This means you don't write a main function for your Encore application.

If you want to customize the initialization behavior of your service, you can define a service struct and define custom initialization logic with that. Read the docs to learn more.