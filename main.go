// (setenv "GOPATH" "/Users/amanda/Bazel/go/.external:/Users/amanda/Bazel/go")

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
	InputType   string
	Input       message
	PathArgs    []field
	Path        string
	Method      string
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
	sources := map[string]string{}

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
			sources["."+file.GetPackage()+"."+msg.GetName()] = file.GetName()
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
					InputType:   meth.GetInputType(),
				}

				if meth.GetOptions() == nil {
					continue
				}

				if tmp, err := proto.GetExtension(meth.GetOptions(), google_api.E_Http); err == nil {
					http := tmp.(*google_api.HttpRule)

					if http.Get != "" {
						m.PathArgs = parsePath(messages, meth, http.Get)
						m.Path = http.Get
						m.Method = "GET"
					}

					if http.Put != "" {
						m.PathArgs = parsePath(messages, meth, http.Put)
						m.Path = http.Put

						m.Method = "PUT"
					}

					if http.Post != "" {
						m.PathArgs = parsePath(messages, meth, http.Post)
						m.Path = http.Post

						m.Method = "POST"
					}

					if http.Delete != "" {
						m.PathArgs = parsePath(messages, meth, http.Delete)
						m.Path = http.Delete

						m.Method = "DELETE"
					}
				}

				s.Methods = append(s.Methods, m)
			}

			if len(s.Methods) > 0 {
				services[svc.GetName()] = s
			}
		}

		if len(services) > 0 {
			fname := strings.Replace(file.GetName(), ".proto", ".pb.kit.go", 1)
			buff := bytes.NewBuffer(nil)
			imps := map[string]string{}
			impUsages := map[string]int{}

			for _, svc := range services {
				for _, meth := range svc.Methods {
					src := sources[meth.InputType]
					if src == file.GetName() {
						continue
					}

					impUsages[src] = impUsages[src] + 1
				}
			}

			for k := range impUsages {
				imps[k] = imports[k]
			}

			fileTemplate.Execute(buff, struct {
				Package  string
				Imports  map[string]string
				Services map[string]service
			}{goPackage, imps, services})
			data := buff.String()

			ret.File = append(ret.File, &plugin.CodeGeneratorResponse_File{
				Name:    &fname,
				Content: &data,
			})
		}
	}
}
