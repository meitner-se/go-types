package types

import (
	"encoding/json"
	"fmt"
)

func ExampleUnmarshal() {
	type Person struct {
		FirstName String
		LastName  String
		Address   String
	}

	jsonString := `{"FirstName": "John", "LastName": null}`

	var person Person
	_ = json.Unmarshal([]byte(jsonString), &person)

	fmt.Printf("FirstName.IsDefined() == %t\t FirstName.IsNil() == %t\t FirstName.String() == '%s'\n",
		person.FirstName.IsDefined(),
		person.FirstName.IsNil(),
		person.FirstName.String(),
	)

	fmt.Printf("LastName.IsDefined() == %t\t LastName.IsNil() == %t\t LastName.String() == '%s'\n",
		person.LastName.IsDefined(), // defined in the json body
		person.LastName.IsNil(),     // null value in the json body
		person.LastName.String(),    // empty string is the default value for nil or undefined strings
	)

	fmt.Printf("Address.IsDefined() == %t\t Address.IsNil() == %t\t Address.String() == '%s'\n",
		person.Address.IsDefined(), // not defined in the json body
		person.Address.IsNil(),     // nil since it wasn't defined in the json body
		person.Address.String(),    // empty string is the default value for nil or undefined strings
	)

	// output:
	// FirstName.IsDefined() == true	 FirstName.IsNil() == false	 FirstName.String() == 'John'
	// LastName.IsDefined() == true	 LastName.IsNil() == true	 LastName.String() == ''
	// Address.IsDefined() == false	 Address.IsNil() == true	 Address.String() == ''
}

func ExampleMarshal() {
	type Person struct {
		FirstName *String `json:"FirstName,omitempty"`
		LastName  *String `json:"LastName,omitempty"`
		Address   *String `json:"Address,omitempty"`
	}

	person := Person{
		FirstName: NewString("").Ptr(),
		LastName:  NewStringFromPtr(nil).Ptr(),
		Address:   NewStringUndefined().Ptr(),
	}

	jsonBytes, _ := json.Marshal(person)

	fmt.Println(string(jsonBytes))

	// output:
	// {"FirstName":"","LastName":null}
}
