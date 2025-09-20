package AST

func (*StatementAssignment) isNode()      {}
func (*StatementAssignment) isStatement() {}

func (*StatementPrint) isNode()      {}
func (*StatementPrint) isStatement() {}

func (*StatementBlock) isNode()      {}
func (*StatementBlock) isStatement() {}

func (*StatementReturn) isNode()      {}
func (*StatementReturn) isStatement() {}

func (*StatementFor) isNode()      {}
func (*StatementFor) isStatement() {}
