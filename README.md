# Native's Effortless Configuration Language

NECL is a simple, human-friendly, structured configuration language.

The syntax for NECL is inspired by [JSON](https://www.json.org), [Groovy](https://groovy-lang.org/), [NGINX configuration](http://nginx.org/en/docs/beginners_guide.html#conf_structure), and others.

NECL uses a key-value structure along with hierarchy for better readability.

## Motivations

Why not use something already used, such as YAML, JSON, etc.?

We find that most configuration languages have some downsides that are quite annoying if you’re using it for complex applications. There’s quite a gap between programming languages and markup languages. For example, in a programming language, to declare a variable you can simply go `x = “hello”` (with obvious variations from language to language), but JSON requires all definitions to be inside brackets: `“foo”: “bar”`. Although that is very good for interoperability, it’s annoying if you’re using it to create a configuration file. Another great example is YAML. YAML files can get quite confusing if you need to create a large one, one can easily get lost in the indentation of the file, that is far more annoying that JSON’s brackets.

NECL attempts to be a bridge between markup / configuration files, and programming languages. It has a syntax similar to what most programming languages look like, but with all the requirements for a markup language. It is made to be easily written and read.

NECL is build around key-value pairs and a well-defined hierarchy that allows for better readability.

## Syntax

Let’s take the following JSON configuration file for a generic application:

```json
{
    "name": "example",
    "description": "Some description",
    "version": "0.0.1",
    "main": "index.js",
    "//": "This is as close to a comment as you are going to get",
    "keywords": ["example", "config"],
    "scripts": {
        "test": "./test.sh",
        "do_stuff": "./do_stuff.sh"
    },
    "bugs": {
        "url": "https://example.com/bugs"
    },
    "contributors": [{
        "name": "John Doe",
        "email": "johndoe@example.com"
    }, {
        "name": "Ivy Lane",
        "url": "https://example.com/ivylane"
    }],
    "dependencies": {
        "dep1": "^1.0.0",
        "dep2": "3.40",
        "dep3": "6.7"
    }
}
```

The NECL equivalent of this configuration is the following:

```
name = "example"
description = "Some description"
version = "0.0.1"
main = "index.js"
// This is a comment
keywords = ["example", "config"]

scripts {
    test = "./test.sh"
    do_stuff = "./do_stuff.sh"
}

bugs {
    url = "https://example.com/bugs"
}

contributors = [
    {
        name = "John Doe"
        email = "johndoe@example.com"
    },
    {
        name = "Ivy Lane",
        url = "https://example.com/ivylane"
    }
]

dependencies {
    dep1 = "^1.0.0",
    dep2 = "3.40",
    dep3 = "6.7"
}
```

We can convert other languages too, for example, this Kubernetes deployment file (YAML):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

Has the following NECL equivalent:

```
apiVersion = "apps/v1"
kind = "Deployment"
metadata {
    name = "nginx-deployment"
    labels {
        app = "nginx"
    }
    // These labels section can also be written as:
    // labels = [
    //     {
    //         app = "nginx"
    //     }
    // ]
    // It's up to you to choose what's prettier
}
spec {
    replicas = 3
    selector {
        matchLabels {
            app = "nginx"
        }
        // This matchLabels section can also be written as the "labels" one.
    }
    template {
        metadata {
            labels {
                app = "nginx"
            }
            // These labels section can also be written as the one above... You got the point
        }
        spec {
            containers {
                // These containers step can be written as 2 types:
                nginx {
                    image = "nginx:1.14.2"
                    ports {
                        containerPort = 80
                    }
                }
                // Or:
                {
                    name = "nginx"
                    image = "nginx:1.14.2"
                    ports {
                        containerPort = 80
                    }
                }
                // It will only depend on how you choose to interpret it on your application
                // The first example will be accessed as "spec.containers.nginx", 
                // whereas the second one will accessed as "spec.containers.name["nginx"]"
            }
        }
    }
}
```

As you can see, NECL can be used in many ways, with better readability and usage.

NECL also supports expressions:

```
// String interpolations
text = "world"
message = "Hello, ${text}!"

// Arithmetic operations
v1 = 1
v2 = 1
sum = v1 + v2

// Functions
up = upper(message)

// Operators
check_sum = v1 == 1 && sum == 2
```

For more information, check the [syntax spec document](SPEC.md)

## References

- [Why JSON isn’t a Good Configuration Language](https://www.lucidchart.com/techblog/2018/07/16/why-json-isnt-a-good-configuration-language/)
- [Don’t Use JSON as a Configuration File Format. (Unless Absolutely You Have To…)](https://revelry.co/insights/development/json-configuration-file-format/)
- [The state of config file formats: XML vs. YAML vs. JSON vs. HCL](https://octopus.com/blog/state-of-config-file-formats)
