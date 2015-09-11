package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gengo/grpc-gateway/third_party/googleapis/google/api"
	"github.com/golang/protobuf/proto"

	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// type Method struct {
// 	Get, Put, Post, Delete *Operation
// }

// func (m *Method) IsUnique() bool {
// 	n := 0
// 	if m.Get != nil {
// 		n++
// 	}

// 	if m.Post != nil {
// 		n++
// 	}

// 	if m.Put != nil {
// 		n++
// 	}

// 	if m.Delete != nil {
// 		n++
// 	}

// 	return n == 1
// }

// type Operation struct {
// 	Input     string
// 	Output    string
// 	Name      string
// 	PathParts []string
// }

func parsePath(msgs map[string]message, method *google_protobuf.MethodDescriptorProto, path string) (ret []field) {
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if len(p) > 0 && p[0] == '{' && p[len(p)-1] == '}' {
			found := false
			for _, field := range msgs[method.GetInputType()].Fields {
				if field.ProtoName == p[1:len(p)-1] {
					ret = append(ret, field)
					found = true
				}
			}

			if !found {
				ret = append(ret, field{"", ""})
			}
		} else {
			ret = append(ret, field{"", ""})
		}
	}

	return
}

type service struct {
	GoName  string
	Methods []method
}

type method struct {
	GoName      string
	GoInputType string
	Input       message
	PathArgs    []field
}

type message struct {
	Fields []field
}

type field struct {
	ProtoName string
	GoName    string
}

func goise(name string) string {
	if tmp := strings.Split(name, "."); len(tmp) > 0 {
		name = tmp[len(tmp)-1]
	}

	return CamelCase(name)
}

func main() {
	msg := plugin.CodeGeneratorRequest{}
	buff, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := proto.Unmarshal(buff, &msg); err != nil {
		panic(err)
	}

	ret := &plugin.CodeGeneratorResponse{}
	defer func() {
		buff, _ := proto.Marshal(ret)
		os.Stdout.Write(buff)
	}()

	param := msg.GetParameter()
	imports := map[string]string{}

	for _, p := range strings.Split(param, ",") {
		if len(p) == 0 {
			continue
		}

		if p[0] == 'M' {
			parts := strings.Split(p[1:], "=")

			imports[parts[0]] = parts[1]
		}
	}

	messages := map[string]message{}

	for _, file := range msg.GetProtoFile() {
		for _, msg := range file.GetMessageType() {
			m := message{}
			for _, f := range msg.GetField() {
				m.Fields = append(m.Fields, field{
					ProtoName: f.GetName(),
					GoName:    goise(f.GetName()),
				})
			}

			messages["."+file.GetPackage()+"."+msg.GetName()] = m
		}
	}

	for _, file := range msg.GetProtoFile() {
		services := map[string]service{}
		goPackage := "main"
		if file.GetOptions() != nil {
			goPackage = file.GetOptions().GetGoPackage()
		}

		for _, svc := range file.GetService() {
			s := service{
				GoName: goise(svc.GetName()),
			}

			for _, meth := range svc.GetMethod() {
				m := method{
					GoName:      goise(meth.GetName()),
					GoInputType: goise(meth.GetInputType()),
					Input:       messages[meth.GetInputType()],
				}

				if tmp, err := proto.GetExtension(meth.GetOptions(), google_api.E_Http); err == nil {
					http := tmp.(*google_api.HttpRule)

					if http.Get != "" {
						m.PathArgs = parsePath(messages, meth, http.Get)
					} else if http.Put != "" {
						m.PathArgs = parsePath(messages, meth, http.Put)
					} else if http.Post != "" {
						m.PathArgs = parsePath(messages, meth, http.Post)
					} else if http.Delete != "" {
						m.PathArgs = parsePath(messages, meth, http.Delete)
					}
				}

				s.Methods = append(s.Methods, m)
			}

			if len(s.Methods) != 0 {
				services[svc.GetName()] = s
			}
		}

		if len(services) != 0 {
			fname := strings.Replace(file.GetName(), ".proto", ".pb.kit.go", 1)
			buff := bytes.NewBuffer(nil)

			fileTemplate.Execute(buff, struct {
				Package  string
				Imports  map[string]string
				Services map[string]service
			}{goPackage, imports, services})
			data := buff.String()

			ret.File = append(ret.File, &plugin.CodeGeneratorResponse_File{
				Name:    &fname,
				Content: &data,
			})
		}
	}
}
