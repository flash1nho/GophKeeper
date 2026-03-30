package print

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/types/known/structpb"
)

func Result(values []*structpb.Value) {
	for _, secretVal := range values {
		fields := secretVal.GetStructValue().GetFields()

		if id, ok := fields["id"]; ok {
			fmt.Printf("id: %v\n", id.GetNumberValue())
		}

		if data, ok := fields["data"]; ok {
			dataFields := data.GetStructValue().GetFields()
			keys := make([]string, 0, len(dataFields))

			for k := range dataFields {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				v := dataFields[k]

				fmt.Printf("%s: %s\n", k, v.GetStringValue())
			}
		}

		if secretType, ok := fields["type"]; ok {
			fmt.Printf("type: %v\n", secretType.GetStringValue())
		}

		if createdAt, ok := fields["created_at"]; ok {
			fmt.Printf("created_at: %v\n", createdAt.GetStringValue())
		}

		if updatedAt, ok := fields["updated_at"]; ok {
			fmt.Printf("updated_at: %v\n", updatedAt.GetStringValue())
		}

		fmt.Println("---")
	}
}
