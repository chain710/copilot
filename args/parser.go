package args

//go:generate pigeon -o parser_generated.go parser.peg
//go:generate gofmt -w parser_generated.go

func ParseAssignStatementList(input string) (AssignStatementList, error) {
	result, err := Parse("", []byte(input))
	if err != nil {
		return nil, err
	}
	return result.(AssignStatementList), nil
}
