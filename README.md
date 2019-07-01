# DemonDin

![Screenshot with Close Button](https://raw.githubusercontent.com/KellyLSB/demondin/master/screenshots/Screenshot%20from%202019-06-30%2019-34-19.png)

```
# git get github.com/KellyLSB/demondin
# cd $GOPATH/go/src/github.com/KellyLSB/demondin
# dep ensure && yarn
# webpack && go run main.go
```

## Framework Layout

# DemonDin relies on three major integrated components: 

- GORM      (Database Object Relational Mapper)
> ^^ Provides a simple model engine; though bulky and easily replaceable.
> ^^ ==> There are some efficiency issues; may require community contribution.
> ^^ ==> GORM/Preloader *[[#1436](https://github.com/jinzhu/gorm/issues/1436)]
> -- Application Logic Fits Here
> vv ==> Modelgen plugin...?
> vv Generates Models [nearly compliant with GORM]
- gQLgen    (GraphQL Object Serializer).
- GoMacaron (main.go / HTTP Router).

### Editing the GraphQL/Database Model Schema 

`graphql/demondin.graphql`:

gQLgen generates the models in `graphql/model/generated.go`

```
# cd graphql; go run github.com/99designs/gqlgen; cd -;
```

gQlgen offers a modelgen plugin; utilizing this may provide a path
for including the extra model structure tags required for GORM's
database automigration tooling.

Currently I'm migrating the extra tags by hand; 
using `# git diff --cached graphql/model/generated.go` for differentiation.

### Delivering content to the client HTTP

`main.go`:

GoMacaron provides interfacing for setting up the HTTP stack.
Sessions are provided to the graphql stack so that data may be delivered.
Templates are configured for rendering frontend applications.
Extensions of core application capacity is configured here.

`templates/default/shop/items.tmpl` | `GET: /shop`:

Delivers the shopfront following application start.

`jsx/checkout/{main, *}.jsx`:

Renders and orchestrates communication between `/graphql`
with a ReactJS application structure.




HeXXeD
