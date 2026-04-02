package helpers

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/pflag"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type FieldInfo struct {
	Key  string
	Type string
}

func PrintResult(values []*structpb.Value) {
	for _, secretVal := range values {
		fields := secretVal.GetStructValue().GetFields()

		if id, ok := fields["id"]; ok {
			fmt.Printf("id: %.0f\n", id.GetNumberValue())
		}

		if fileName, ok := fields["file_name"]; ok {
			val := fileName.GetStringValue()

			if val != "" {
				fmt.Printf("file_name: %s\n", val)
			}
		}

		if data, ok := fields["data"]; ok && data.GetStructValue() != nil {
			dataFields := data.GetStructValue().GetFields()
			keys := make([]string, 0, len(dataFields))

			for k := range dataFields {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				v := dataFields[k]

				fmt.Printf("%s: %v\n", k, v.AsInterface())
			}
		}

		meta := []string{"type", "created_at", "updated_at"}

		for _, m := range meta {
			if val, ok := fields[m]; ok {
				fmt.Printf("%s: %s\n", m, val.GetStringValue())
			}
		}

		fmt.Println("---")
	}
}

func ArgsParse(cmd *cobra.Command) (int, *structpb.Struct, string, error) {
	dataMap := make(map[string]interface{})

	var id int

	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Name == "id" {
			val, _ := cmd.Flags().GetInt("id")
			id = val

			return
		}

		var val interface{}

		switch f.Value.Type() {
		case "int", "int32", "int64":
			val, _ = cmd.Flags().GetInt64(f.Name)
		case "bool":
			val, _ = cmd.Flags().GetBool(f.Name)
		case "float64", "float32":
			val, _ = cmd.Flags().GetFloat64(f.Name)
		default:
			val = f.Value.String()
		}

		dataMap[f.Name] = val
	})

	data, err := structpb.NewStruct(dataMap)

	if err != nil {
		return 0, nil, "", err
	}

	var secretType string

	if cmd.Parent() != nil {
		secretType = strcase.ToCamel(cmd.Parent().Name())
	}

	return id, data, secretType, nil
}

func ErrorHandler(log *zap.Logger, err error) {
	if statusErr, ok := status.FromError(err); ok {
		fmt.Printf("❌ %s\n", statusErr.Message())
	} else {
		log.Error(err.Error())
	}
}

func GetStructKeys(s interface{}) []FieldInfo {
	var info []FieldInfo
	val := reflect.TypeOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Tag.Get("ignore") == "true" {
			continue
		}

		json := field.Tag.Get("json")

		if json == "" || json == "-" {
			continue
		}

		key := strings.Split(json, ",")[0]

		info = append(info, FieldInfo{
			Key:  key,
			Type: field.Type.String(),
		})
	}

	return info
}
