# CRUDer

Creates the needed golang code to offer a REST endpoint from a struct type definition. 
Generated bits include the listening endpoint, the logic and the database operations. 

You simply need to create a golang struct type in a file and execute this tool. The result
is a full implemented REST service based on that struct type.

## Install

### From debian package

Add private ppa:

```sh

```
and install

```sh
sudo apt update
sudo apt install cruder
```

### From snap

```sh
snap install cruder
```

## How to create a service

Just move to the project folder

```sh
cd $GOPATH/src/theserver.com/youruser/yourproject
```

Open your favourite editor and edit a new file where defining your type:

```sh
vi mytype.go
```

and populate its content with something like:

```golang
package whatever

type MyType struct {
        ID          int
        Name        string
        Description string
        Whatever    bool
}
```

save the content and launch cruder tool:

```sh
cruder mytype.go
```

At this point the service is created. Let's provide a settings file to test it works:

```sh
$ cat <<EOF > settings.yaml
host: localhost
port: 8080
driver: sqlite3
datasource: ./main.db

EOF
```

and launch:

```sh
$ go run cmd/service/main.db
2018/04/02 15:42:19 Started service on port 8080
```

There you go!. Open your browser and type:

```web
http://localhost:8080/v1/mytype
```

and you should get a reply like `{"mytypes":[]}`

You can also use _curl_ from command line:

```sh
$ curl -i -X GET http://localhost:8080/v1/mytype
HTTP/1.1 200 OK
Date: Mon, 02 Apr 2018 13:52:31 GMT
Content-Length: 15
Content-Type: text/plain; charset=utf-8

{"mytypes":[]}
```

## What has been created?

You can check the generated files and folders by showing the tree 
in the root of the project:

```sh
$ tree -d
.
├── cmd
│   └── service
│       └── main.go
├── datastore
│   ├── db.go
│   ├── ddl.go
│   └── mytype.go
├── handler
│   ├── mytype.go
│   └── reply.go
├── mytype.go
└── service
    ├── router.go
    └── service.go
```

As you can see, there are several subfolders and created files:

- cmd/service/main.go file holds the entry point to the service
- datastore folder includes all the operational bits to access database
  - _db.go_: generic database definition and opening
  - _ddl.go_: data definition language operations, including tables creation, alter, etc..
  - _mytype.go_: database operations related to just created type. The name of this file
  is the name of the provided type and the file itself includes the provided type definition.
- handler folder holds the REST logic layer
  - _mytype.go_: includes REST endpoint operations related with provided type. The
  name of this file depends on the name of the provided type.
  - _reply.go_: generic response helper methods
- service folder includes general service files
  - _router.go_: includes all the exposed routes of REST operations. It has new entries for the
  CRUD operations for the provided type
  - _service.go_: bits to start the service with a provided configuration

These are the default generated files by cruder, but you can modify its behaviour
by updating/replacing/removing the built-in plugins that generates this data.

## Additional types

Project has not to be necessarily empty before CRUDer execution. You can add the additional types that you need to the REST service. Every CRUDer execution generates the new type related files, preserving and modifying the previous ones for all them to be part of the same service.

You can define another type like:

```golang
package anyother

type AnotherType struct {
    ID          int
    Name        string
    Description string
    Exists      bool
}
```

and execute cruder command again:

```sh
cruder anothertype.go
```

resulting in the service completed like:

```sh
$ tree 
.
├── anothertype.go
├── cmd
│   └── service
│       └── main.go
├── datastore
│   ├── anothertype.go
│   ├── db.go
│   ├── ddl.go
│   └── mytype.go
├── handler
│   ├── anothertype.go
│   ├── mytype.go
│   └── reply.go
├── main.db
├── mytype.go
├── service
│   ├── router.go
│   └── service.go
└── settings.yaml
```

## Plugins

Generated code is defined by a set of plugins, some of them included in CRUDer distribution, some
defined by the user.

### Built-it

All the default generated code is created by some plugins that are distributed along with CRUDer.
You can find them under `/usr/lib/cruder/plugins/` as `.so` shared library files.

- _datastore.so_ plugin generates `datastore/mytype.go` file
- _db.so_ plugin generates `datastore/db.go`file
- _ddl.so_ plugin generates `datastore/ddl.go`file
- _handler.so_ plugin generates `handler/mytype.so` file
- _main.so_ plugin generates `cmd/service/main.go` file
- _reply.so_ plugin generates `handler/reply.go` file
- _router.so_ plugin generates `service/router.go` file
- _service.so_ plugin generates `service/service.go` file

### User defined

You can get rid of not desired built-in plugins, overwrite them or create additional ones. Building you own
plugin is really easy. In general, you just have to write a template and a golang file.

## Plugin development

If you desire to build your own plugins you can do it by writing a template + transformation code.

### Template

The first step is defining your template. A template is simply a golang source file populated with literals and placeholders. Literal content is not modified, and placeholders are replaced by specific values of the defined type when cruder is executed. The name of the template must be the identifier of the plugin with .template extension. Like:

```sh
myplugin.template
```

A fragment of the template could be similar to:

```golang
type Datastore interface {
  Create_#TYPE#_Table() error
  List_#TYPE#_s() ([]_#TYPE#_, error)
  Get_#TYPE#_(_#ID.FIELD.NAME#_ _#ID.FIELD.TYPE#_) (_#TYPE#_, error)
  Find_#TYPE#_(query string) (_#TYPE#_, error)
  Create_#TYPE#_(_#TYPE.IDENTIFIER#_ _#TYPE#_) (int, error)
  Update_#TYPE#_(_#ID.FIELD.NAME#_ _#ID.FIELD.TYPE#_, _#TYPE.IDENTIFIER#_ _#TYPE#_)
  Delete_#TYPE#_(_#ID.FIELD.NAME#_ _#ID.FIELD.TYPE#_) error
}
```

The placeholders are the different symbols between marks `_#` and `#_`. They are the fixed values defined in this table:

| Placeholder | TheType value | Description |
| ----------- | :------------ | :---------- |
| \_#TYPE#\_ | TheType | Name of the type |
| \_#TYPE.IDENTIFIER#\_ | theType | Type identifier|
| \_#TYPE.LOWERCASE#\_ | thetype | Type in lower case |
| \_#ID.FIELD.NAME#\_ | ID | Identifier field name |
| \_#ID.FIELD.NAME.LOWERCASE#\_ | id | Identifier field of the type in lower case |
| \_#ID.FIELD.TYPE#\_ | int | Identifier field type |
| \_#FIND.FIELD.NAME#\_ | name | Name of the field used for searching |
| \_#FIELDS.ENUM#\_ | theType.Name, theType.Description, theType.Subtypes | Enumeration of type fields |
| \_#FIELDS.ENUM.REF#\_ | &theType.Name, &theType.Description, &theType.Subtypes | Enumeration of reference type fields |
| \_#ID.FIELD.DDL#\_ | id integer primary_key not null | Identifier field in DDL sentences |
| \_#FIELDS.DDL#\_ | Name	varchar,\nDescription	varchar,\nSubtypes varchar | Fields in DDL sentences |
| \_#FIELDS.DML#\_ | name, description, subtypes | Fields in DML sentences |
| \_#VALUES.DML.PARAMS#\_ | $1, $2, $3 | Params in DML sentences |
| \_#ID.FIELD.DML.PARAM#\_ | id=$4 | Identifier field in DML sentences |
| \_#FIELDS.DML.PARAMS#\_ | name=$1, description=$2, subtypes=$3 | Fields in DML sentences |
| \_#ID.FIELD.TYPE.PARSE#\_ | strconv.Atoi(vars["id"]) | Conversion instruction for identifier field from string to its type |
| \_#ID.FIELD.TYPE.FORMAT#\_ | strconv.Itoa(id) | Conversion instruction for identifier field from its type to string |
| \_#ID.FIELD.PATTERN#\_ | [a-z]+ | Regular expression matching possible identifier field values |

NOTE: consider *TheType* like:

```golang
type TheType struct {
  ID          int
  Name        string
  Description string
  SubTypes    bool
}
```

This would generate the output:

```golang
type Datastore interface {
  CreateTheTypeTable() error
  ListTheTypes() ([]TheType, error)
  GetTheType(id int) (TheType, error)
  FindTheType(name string) (TheType, error)
  CreateTheType(theType TheType) (int, error)
  UpdateTheType(id int, theType TheType)
  DeleteTheType(id int) error
}
```

There are additional placeholders related with a general configuration:

| Placeholder | sample value | Description |
| ----------- | :------------ | :---------- |
| _#PROJECT#_ | github.com/myuser/myproject | import path for current project |
| _#API.VERSION#_ | v1 | Version of the exposed API |


### Transformation code

This is the code that will operate over the template and type merge output. In general, you just have to implement `makers.Maker` interface, that defines the desired behaviour of the plugin when replaced template placeholders with the provided type related values.

```golang
type Maker interface {
  ID() string
  Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error)
  OutputFilepath() string
}
```

#### Development

1.- To start developing your own transformation code, you have to create a new file and make it belong to `main` package:

```golang
package main
...
```

2.- Below the package definition, define the struct that is going to be your plugin. For an easier development, CRUDer provides an anonymous base maker that you should include in your struct, like this:

```golang
package main

import (
  "github.com/rmescandon/cruder/makers"
)

type MyPlugin struct {
  makers.Base
}
```

3.- Now, time to implement the methods of `makers.Maker` interface. Let's start with returning an identifier for the plugin. This shouldn't match any of the existing plugins, built-in included. So, take care of not selecting *ddl*, *handler*, *main*, *reply*, *router*, *service*, *db*, *datastore* or any other plugin identifier you have added before.

```golang
func (p *MyPlugin) ID() string {
  return "myplugin"
}
```

4.- Implement the method that says CRUDer where the resultant file must be written and its name:

```golang
func (p *MyPlugin) OutputFilepath() string {
  return filepath.Join(makers.BasePath, "mypluginfolder", p.ID()+".go")
}
```

`makers.BasePath` points to the output folder selected as param when cruder is executed. Default is current directory.

5.- Implement the method that defines the specific transformation for this plugin

```golang
func (p *MyPlugin) Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error) {
  ...
}
```

Here are the params explanation:

- generatedOutput: Is the content that CRUDer creates by merging template with provided type.
- currentOutput: Is the content of a generated file by this plugin in a previous execution of CRUDer

The returned content should have what must be written to output file. If null is returned, nothing new is written to output (if a previous file existed, it is not overwriten). A returned error won't stop processing the rest of the plugins

6.- Finally, register your plugin in the init() function. This lets CRUDer engine include your plugin in the list of available ones.

```golang
func init() {
	makers.Register(&MyPlugin{})
}
```

#### Compile

You need to compile the just created plugin transformation code by executing:
```sh
go build -buildmode=plugin myplugin.go
```

This generates `myplugin.so` as output file

#### Deploy

Move `myplugin.template` to `/usr/share/cruder/templates` and `myplugin.so` to `/usr/lib/cruder/plugins`

That's it. The plugin will be used in next CRUDer execution

#### A sample transformation code

For example, the Reply plugin, that only makes a copy of the generated output, has this code:

```golang
package main

import (
  "path/filepath"

  "github.com/rmescandon/cruder/errs"
  "github.com/rmescandon/cruder/io"
  "github.com/rmescandon/cruder/makers"
)

// Reply struct holding data to copy reply template
type Reply struct {
  makers.Base
}

// ID returns 'reply' as this maker identifier
func (r *Reply) ID() string {
  return "reply"
}

// OutputFilepath returns the path to the output file
func (r *Reply) OutputFilepath() string {
  return filepath.Join(makers.BasePath, "handler/reply.go")
}

// Make copies template to output path
func (r *Reply) Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error) {
  if currentOutput != nil {
    return nil, errs.NewErrOutputExists(r.OutputFilepath())
  }

  return generatedOutput, nil
}

func init() {
  makers.Register(&Reply{})
}
```

You can find all the build-in plugins source code under `makers/plugins`