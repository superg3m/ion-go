package AST

func (*DeclarationVariable) isNode()        {}
func (*DeclarationVariable) isDeclaration() {}

func (*DeclarationFunction) isNode()        {}
func (*DeclarationFunction) isDeclaration() {}

func (d DeclarationStruct) isNode()        {}
func (d DeclarationStruct) isDeclaration() {}
