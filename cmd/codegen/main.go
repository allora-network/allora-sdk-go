package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	butils "github.com/brynbellomy/go-utils"
	annotations "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

//go:embed client.wrapper.go.tpl
var clientWrapperTemplate string

//go:embed client.wrapper.module.go.tpl
var clientWrapperModuleTemplate string

//go:embed client.rest.module.go.tpl
var clientRestModuleTemplate string

//go:embed client.grpc.go.tpl
var clientGRPCTemplate string

//go:embed client.grpc.module.go.tpl
var clientGRPCModuleTemplate string

//go:embed client.rest.go.tpl
var restAggregatedClientTemplate string

//go:embed client.interface.module.go.tpl
var clientInterfaceModuleTemplate string

//go:embed client.interface.go.tpl
var clientInterfaceTemplate string

type ClientMethod struct {
	Name         string
	RequestType  string
	ResponseType string
	Comment      string
	// HTTP-specific fields for REST clients
	HTTPMethod  string   // GET, POST, etc.
	HTTPPath    string   // /cosmos/bank/v1beta1/balances/{address}
	PathParams  []string // ["address"]
	QueryParams []string // ["denom", "pagination.limit"]
}

type ModuleInfo struct {
	ModuleName  string
	PackageName string
	ImportPath  string
	ServiceName string
	Methods     []ClientMethod
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/codegen/main.go <command> [args...]\n" +
			"Commands:\n" +
			"  generate [flags]       Generate wrappers from .proto files\n\n" +
			"Flags for generate:\n" +
			"  -I <path>   Include directory (repeatable)\n" +
			"  -f <file>   Proto file to process (repeatable)\n")
	}

	generateFromProtos(os.Args[1:])
}

// stringSlice supports repeatable flags
type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func generateFromProtos(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	var includes stringSlice
	var files stringSlice
	fs.Var(&includes, "I", "include path")
	fs.Var(&files, "f", "proto file")
	_ = fs.Parse(args)
	if len(files) == 0 {
		log.Fatal("no proto files provided; pass at least one -f <file>")
	}

	// Build descriptor set from all provided files
	fds, err := buildFileDescriptorSetForFiles(files, includes)
	if err != nil {
		log.Fatalf("failed to build descriptor set: %v", err)
	}

	// Track generated files for post-generation formatting
	generatedFiles := []string{}

	// Map proto package -> go_package import path and optional name override
	type goPkgInfo struct{ importPath, pkgName string }
	pkgToGo := map[string]goPkgInfo{}
	for _, fp := range fds.File {
		pkg := fp.GetPackage()
		goPkg := fp.GetOptions().GetGoPackage()
		if pkg == "" || goPkg == "" {
			continue
		}
		imp := goPkg
		name := ""
		if semi := strings.Index(goPkg, ";"); semi >= 0 {
			imp = goPkg[:semi]
			name = goPkg[semi+1:]
		}

		switch {
		case strings.HasPrefix(imp, "github.com/cosmos/cosmos-sdk/x/epochs"):
			imp = "cosmossdk.io/x/epochs" + imp[len("github.com/cosmos/cosmos-sdk/x/epochs"):]
		case strings.HasPrefix(imp, "github.com/cosmos/cosmos-sdk/x/evidence"):
			imp = "cosmossdk.io/x/evidence" + imp[len("github.com/cosmos/cosmos-sdk/x/evidence"):]
		case strings.HasPrefix(imp, "github.com/cosmos/cosmos-sdk/x/feegrant"):
			imp = "cosmossdk.io/x/feegrant" + imp[len("github.com/cosmos/cosmos-sdk/x/feegrant"):]
		case strings.HasPrefix(imp, "github.com/cosmos/cosmos-sdk/x/protocolpool"):
			imp = "cosmossdk.io/x/protocolpool" + imp[len("github.com/cosmos/cosmos-sdk/x/protocolpool"):]
		}

		pkgToGo[pkg] = goPkgInfo{importPath: imp, pkgName: name}
	}

	filesReg, err := protodesc.NewFiles(fds)
	if err != nil {
		log.Fatalf("failed to create files registry: %v", err)
	}

	fileSet := butils.NewSet[string]()
	for _, file := range files {
		for _, include := range includes {
			if !strings.HasPrefix(file, include) {
				continue
			}

			relPath, err := filepath.Rel(include, file)
			if err != nil {
				log.Fatalf("failed to compute relative path for %s (include %s): %v", file, include, err)
			}
			fileSet.Add(relPath)
		}
	}

	var modules []ModuleInfo
	filesReg.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if !fileSet.Has(fd.Path()) {
			return true
		}

		protoPkg := string(fd.Package())
		info, ok := pkgToGo[protoPkg]
		if !ok || info.importPath == "" {
			return true
		}

		moduleName := inferModuleNameFromProtoPackage(protoPkg)
		pkgAlias := computeImportAlias(info.importPath, moduleName, info.pkgName)

		for i := 0; i < fd.Services().Len(); i++ {
			sd := fd.Services().Get(i)
			methods := parseMethodsForService(sd)
			if len(methods) == 0 {
				continue
			}
			modules = append(modules, ModuleInfo{
				ModuleName:  moduleName,
				PackageName: pkgAlias,
				ImportPath:  info.importPath,
				ServiceName: string(sd.Name()),
				Methods:     methods,
			})
		}
		return true
	})

	if len(modules) == 0 {
		log.Println("No modules discovered in provided protos")
		return
	}
	for _, m := range modules {
		fmt.Printf("Generating %s module...\n", m.ModuleName)

		name := fmt.Sprintf("gen/wrapper/client.%s.go", m.ModuleName)
		if err := executeTemplate(m, name, clientWrapperModuleTemplate); err != nil {
			log.Fatalf("failed to generate wrapper client: %v", err)
		}
		generatedFiles = append(generatedFiles, name)

		name = fmt.Sprintf("gen/rest/client.%s.go", m.ModuleName)
		if err := executeTemplate(m, name, clientRestModuleTemplate); err != nil {
			log.Fatalf("failed to generate REST client: %v", err)
		}
		generatedFiles = append(generatedFiles, name)

		name = fmt.Sprintf("gen/grpc/client.%s.go", m.ModuleName)
		if err := executeTemplate(m, name, clientGRPCModuleTemplate); err != nil {
			log.Fatalf("failed to generate gRPC client: %v", err)
		}
		generatedFiles = append(generatedFiles, name)

		name = fmt.Sprintf("gen/interfaces/client.%s.go", m.ModuleName)
		if err := executeTemplate(m, name, clientInterfaceModuleTemplate); err != nil {
			log.Fatalf("failed to generate module client interface (%v): %v", m.ModuleName, err)
		}
		generatedFiles = append(generatedFiles, name)
	}

	data := struct {
		Modules []ModuleInfo
	}{
		modules,
	}

	name := "gen/interfaces/client.go"
	if err := executeTemplate(data, name, clientInterfaceTemplate); err != nil {
		log.Fatalf("failed to generate aggregated client: %v", err)
	}
	generatedFiles = append(generatedFiles, name)

	// Generate REST aggregated client
	name = "gen/rest/client.go"
	if err := executeTemplate(data, name, restAggregatedClientTemplate); err != nil {
		log.Fatalf("failed to generate REST aggregated client: %v", err)
	}
	generatedFiles = append(generatedFiles, name)

	// Generate gRPC aggregated client
	name = "gen/grpc/client.go"
	if err := executeTemplate(data, name, clientGRPCTemplate); err != nil {
		log.Fatalf("failed to generate gRPC aggregated client: %v", err)
	}
	generatedFiles = append(generatedFiles, name)

	name = "gen/wrapper/client.go"
	if err := executeTemplate(data, name, clientWrapperTemplate); err != nil {
		log.Fatalf("failed to generate gRPC aggregated client: %v", err)
	}
	generatedFiles = append(generatedFiles, name)

	// Format generated files with gofmt and goimports (if available)
	if err := formatGeneratedFiles(generatedFiles); err != nil {
		log.Printf("warning: failed to format generated files: %v", err)
	}

	fmt.Println("Done generating wrappers from protos.")
}

func inferModuleNameFromProtoPackage(protoPkg string) string {
	parts := strings.Split(protoPkg, ".")
	end := len(parts)
	for end > 0 && isVersionToken(parts[end-1]) {
		end--
	}
	if end == 0 {
		return "module"
	}
	return parts[end-1]
}

func isVersionToken(tok string) bool {
	return regexp.MustCompile(`^v[0-9]+((alpha|beta)[0-9]*)?`).MatchString(tok)
}

func computeImportAlias(importPath, moduleName, pkgName string) string {
	if pkgName != "" {
		return pkgName
	}
	base := path.Base(importPath)
	switch base {
	case "types", "v1", "v1beta1", "v2", "v2beta1", "v3":
		return moduleName + base
	default:
		return base
	}
}

func buildFileDescriptorSetForFiles(files []string, includeDirs []string) (*descriptorpb.FileDescriptorSet, error) {
	tmpFile := filepath.Join(os.TempDir(), "descriptor-set.pb")
	args := []string{}
	for _, inc := range includeDirs {
		args = append(args, "-I", inc)
	}
	args = append(args, "--include_imports", "--include_source_info", "--descriptor_set_out="+tmpFile)
	args = append(args, files...)

	cmd := exec.Command("protoc", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("protoc failed: %v, output: %s", err, string(out))
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read descriptor set: %w", err)
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal descriptor set: %w", err)
	}
	return fds, nil
}

func parseMethodsForService(sd protoreflect.ServiceDescriptor) []ClientMethod {
	var methods []ClientMethod
	for j := 0; j < sd.Methods().Len(); j++ {
		md := sd.Methods().Get(j)
		cm := ClientMethod{
			Name:         string(md.Name()),
			RequestType:  string(md.Input().Name()),
			ResponseType: string(md.Output().Name()),
		}

		if opts := md.Options(); opts != nil {
			if v := proto.GetExtension(opts.(*descriptorpb.MethodOptions), annotations.E_Http); v != nil {
				if httpRule, ok := v.(*annotations.HttpRule); ok && httpRule != nil {
					method, path := httpRulePattern(httpRule)
					cm.HTTPMethod = method
					cm.HTTPPath = sanitizeHTTPPath(path)
					extractPathParams(&cm)
				}
			}
		}
		// derive query params using only the method descriptor context
		deriveQueryParams(&cm, md)
		methods = append(methods, cm)
	}
	return methods
}

func httpRulePattern(rule *annotations.HttpRule) (string, string) {
	if p := rule.GetGet(); p != "" {
		return "GET", p
	} else if p := rule.GetPost(); p != "" {
		return "POST", p
	} else if p := rule.GetPut(); p != "" {
		return "PUT", p
	} else if p := rule.GetDelete(); p != "" {
		return "DELETE", p
	} else if p := rule.GetPatch(); p != "" {
		return "PATCH", p
	} else if c := rule.GetCustom(); c != nil {
		return c.Kind, c.Path
	}
	return "", ""
}

// sanitizeHTTPPath removes wildcard patterns inside path params, e.g. {denom=**} -> {denom}
func sanitizeHTTPPath(path string) string {
	// Pattern explanation:
	// \{        - literal opening brace
	// ([^}=]+)  - capture group 1: one or more chars that are not '}' or '='
	// (?:=.*?)?  - non-capturing group: optional '=' followed by any chars (non-greedy)
	// \}        - literal closing brace
	re := regexp.MustCompile(`\{([^}=]+)(?:=.*?)?\}`)
	return re.ReplaceAllString(path, "{$1}")
}

func extractPathParams(method *ClientMethod) {
	path := method.HTTPPath
	for i := 0; i < len(path); i++ {
		if path[i] == '{' {
			j := i + 1
			for j < len(path) && path[j] != '}' {
				j++
			}
			if j < len(path) {
				param := path[i+1 : j]
				if eq := strings.IndexByte(param, '='); eq >= 0 {
					param = param[:eq]
				}
				if param != "" {
					method.PathParams = append(method.PathParams, param)
				}
				i = j
			}
		}
	}
}

func deriveQueryParams(cm *ClientMethod, md protoreflect.MethodDescriptor) {
	if cm.HTTPMethod == "" {
		return
	}
	req := md.Input()
	pathParamSet := map[string]struct{}{}
	for _, p := range cm.PathParams {
		pathParamSet[p] = struct{}{}
	}

	// If a 'pagination' field exists, add standard pagination params
	foundPagination := false
	for i := 0; i < req.Fields().Len(); i++ {
		f := req.Fields().Get(i)
		if string(f.Name()) == "pagination" {
			foundPagination = true
			break
		}
	}
	if foundPagination {
		cm.QueryParams = append(cm.QueryParams,
			"pagination.key",
			"pagination.offset",
			"pagination.limit",
			"pagination.count_total",
			"pagination.reverse",
		)
	}

	// Add simple scalars not in path
	for i := 0; i < req.Fields().Len(); i++ {
		f := req.Fields().Get(i)
		name := string(f.Name())
		if name == "pagination" {
			continue
		}
		if _, isPath := pathParamSet[name]; isPath {
			continue
		}
		if f.Kind() != protoreflect.MessageKind {
			cm.QueryParams = append(cm.QueryParams, name)
		}
	}
}

func executeTemplate(this any, outfile string, tmpl string) error {
	// Parse template with custom functions
	funcMap := template.FuncMap{
		"title": strings.Title,
	}

	t, err := template.New(outfile).Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse %s template: %w", outfile, err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outfile), 0o755); err != nil {
		return fmt.Errorf("failed to create output dir %s: %w", filepath.Dir(outfile), err)
	}

	// Create output file
	file, err := os.Create(outfile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outfile, err)
	}
	defer file.Close()

	// Execute template
	if err := t.Execute(file, this); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", outfile, err)
	}

	fmt.Printf("Generated %s\n", outfile)
	return nil
}

// formatGeneratedFiles runs gofmt and goimports over the provided files.
// If goimports is not installed, it is skipped with a warning.
func formatGeneratedFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}

	fmt.Printf("Formatting %d generated files...\n", len(files))

	// Run: gofmt -w <files...>
	gofmtArgs := append([]string{"-w"}, files...)
	if out, err := exec.Command("gofmt", gofmtArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("gofmt failed: %v, output: %s", err, string(out))
	}

	// Run: goimports -w <files...> (if available)
	if _, err := exec.LookPath("goimports"); err == nil {
		goimportsArgs := append([]string{"-w"}, files...)
		if out, err := exec.Command("goimports", goimportsArgs...).CombinedOutput(); err != nil {
			return fmt.Errorf("goimports failed: %v, output: %s", err, string(out))
		}
	} else {
		log.Printf("goimports not found in PATH; skipping goimports. Install with: go install golang.org/x/tools/cmd/goimports@latest")
	}

	return nil
}
