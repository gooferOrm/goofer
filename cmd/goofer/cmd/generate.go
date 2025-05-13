package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	entityName   string
	fields       []string
	outputDir    string
	packageName  string
	withValidate bool
	withHooks    bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code for Goofer ORM entities",
	Long: `Generate boilerplate code for Goofer ORM entities, DTOs, repositories, and more.
This command helps you quickly scaffold the code needed to work with Goofer ORM.`,
}

// entityCmd represents the entity generate command
var entityCmd = &cobra.Command{
	Use:   "entity [name] [field:type:tag...]",
	Short: "Generate an entity struct",
	Long: `Generate a new entity struct with the specified fields and ORM tags.

Example:
  goofer generate entity User id:uint:primaryKey,autoIncrement name:string:notnull email:string:unique,notnull

Field types: string, int, uint, float64, bool, time.Time
Tags: primaryKey, autoIncrement, unique, notnull, index`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityName = args[0]
		if len(args) > 1 {
			fields = args[1:]
		}
		generateEntity()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(entityCmd)

	// Add flags for entity command
	entityCmd.Flags().StringVarP(&outputDir, "out", "o", ".", "Output directory for generated code")
	entityCmd.Flags().StringVarP(&packageName, "package", "p", "models", "Package name for generated code")
	entityCmd.Flags().BoolVar(&withValidate, "with-validate", false, "Add validation tags")
	entityCmd.Flags().BoolVar(&withHooks, "with-hooks", false, "Add lifecycle hooks")
}

func generateEntity() {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Parse field definitions
	var parsedFields []FieldDefinition
	for _, field := range fields {
		fieldDef := parseFieldDefinition(field)
		parsedFields = append(parsedFields, fieldDef)
	}

	// Add default ID field if not present
	hasID := false
	for _, field := range parsedFields {
		if field.Name == "ID" {
			hasID = true
			break
		}
	}

	if !hasID {
		parsedFields = append([]FieldDefinition{
			{
				Name:    "ID",
				Type:    "uint",
				OrmTags: []string{"primaryKey", "autoIncrement"},
			},
		}, parsedFields...)
	}

	// Add timestamps if with-hooks is enabled
	if withHooks {
		hasCreatedAt := false
		hasUpdatedAt := false

		for _, field := range parsedFields {
			if field.Name == "CreatedAt" {
				hasCreatedAt = true
			}
			if field.Name == "UpdatedAt" {
				hasUpdatedAt = true
			}
		}

		if !hasCreatedAt {
			parsedFields = append(parsedFields, FieldDefinition{
				Name:    "CreatedAt",
				Type:    "time.Time",
				OrmTags: []string{"type:timestamp", "default:CURRENT_TIMESTAMP"},
			})
		}

		if !hasUpdatedAt {
			parsedFields = append(parsedFields, FieldDefinition{
				Name:    "UpdatedAt",
				Type:    "time.Time",
				OrmTags: []string{"type:timestamp"},
			})
		}
	}

	// Prepare the template data
	data := EntityTemplateData{
		PackageName: packageName,
		EntityName:  entityName,
		Fields:      parsedFields,
		WithHooks:   withHooks,
	}

	// Generate the entity code
	filePath := filepath.Join(outputDir, strings.ToLower(entityName)+".go")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Execute the template
	err = entityTemplate.Execute(file, data)
	if err != nil {
		fmt.Printf("Error generating entity: %v\n", err)
		return
	}

	fmt.Printf("Generated entity %s in %s\n", entityName, filePath)
}

// FieldDefinition represents a field in an entity
type FieldDefinition struct {
	Name       string
	Type       string
	OrmTags    []string
	ValidTags  []string
	IsRequired bool
}

// EntityTemplateData contains data for entity template
type EntityTemplateData struct {
	PackageName string
	EntityName  string
	Fields      []FieldDefinition
	WithHooks   bool
}

// parseFieldDefinition parses a field definition string
func parseFieldDefinition(fieldDef string) FieldDefinition {
	parts := strings.Split(fieldDef, ":")
	name := parts[0]
	fieldType := "string" // default type
	var ormTags []string
	var validTags []string
	isRequired := false

	if len(parts) > 1 {
		fieldType = parts[1]
	}

	if len(parts) > 2 {
		tagStr := parts[2]
		tags := strings.Split(tagStr, ",")
		
		for _, tag := range tags {
			ormTags = append(ormTags, tag)
			
			// Add corresponding validation tags
			if withValidate {
				if tag == "notnull" {
					validTags = append(validTags, "required")
					isRequired = true
				}
				
				if fieldType == "string" {
					validTags = append(validTags, "max=255")
				}
				
				if tag == "unique" && fieldType == "string" && name == "Email" {
					validTags = append(validTags, "email")
				}
			}
		}
	}

	return FieldDefinition{
		Name:       name,
		Type:       fieldType,
		OrmTags:    ormTags,
		ValidTags:  validTags,
		IsRequired: isRequired,
	}
}

// Format field tags for the template
func (f FieldDefinition) FormatOrmTags() string {
	if len(f.OrmTags) == 0 {
		return ""
	}
	return fmt.Sprintf(`orm:"%s"`, strings.Join(f.OrmTags, ";"))
}

func (f FieldDefinition) FormatValidateTags() string {
	if !withValidate || len(f.ValidTags) == 0 {
		return ""
	}
	return fmt.Sprintf(`validate:"%s"`, strings.Join(f.ValidTags, ","))
}

func (f FieldDefinition) FormatTags() string {
	ormTag := f.FormatOrmTags()
	validateTag := f.FormatValidateTags()
	
	if ormTag != "" && validateTag != "" {
		return fmt.Sprintf("`%s %s`", ormTag, validateTag)
	} else if ormTag != "" {
		return fmt.Sprintf("`%s`", ormTag)
	} else if validateTag != "" {
		return fmt.Sprintf("`%s`", validateTag)
	}
	
	return ""
}

// toLowerCase is a helper function for the template
func toLowerCase(s string) string {
	return strings.ToLower(s)
}

// Template for entity generation
var entityTemplate *template.Template

func init() {
	// Create a new template with our custom functions
	entityTemplate = template.New("entity").Funcs(template.FuncMap{
		"toLowerCase": toLowerCase,
	})

	// Parse the template
	template.Must(entityTemplate.Parse(`package {{ .PackageName }}

import (
{{- if .WithHooks }}
	"time"
	"fmt"
{{- else }}
	"time"
{{- end }}

	"github.com/gooferOrm/goofer/schema"
)

// {{ .EntityName }} entity
type {{ .EntityName }} struct {
{{- range .Fields }}
	{{ .Name }} {{ .Type }} {{ .FormatTags }}
{{- end }}
}

// TableName returns the table name for the {{ .EntityName }} entity
func ({{ .EntityName }}) TableName() string {
	return "{{ .EntityName | toLowerCase }}s"
}
{{ if .WithHooks }}

// BeforeCreate is called before creating a new record
func (e *{{ .EntityName }}) BeforeCreate() error {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is called before updating a record
func (e *{{ .EntityName }}) BeforeUpdate() error {
	e.UpdatedAt = time.Now()
	return nil
}

// AfterCreate is called after creating a new record
func (e *{{ .EntityName }}) AfterCreate() error {
	fmt.Printf("{{ .EntityName }} created with ID: %v\n", e.ID)
	return nil
}

// AfterUpdate is called after updating a record
func (e *{{ .EntityName }}) AfterUpdate() error {
	fmt.Printf("{{ .EntityName }} updated: %v\n", e.ID)
	return nil
}

// BeforeDelete is called before deleting a record
func (e *{{ .EntityName }}) BeforeDelete() error {
	fmt.Printf("About to delete {{ .EntityName }} with ID: %v\n", e.ID)
	return nil
}
{{ end }}
`))
}