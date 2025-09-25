package AST

func (*StatementAssignment) isNode()      {}
func (*StatementAssignment) isStatement() {}

func (*StatementPrint) isNode()       {}
func (*StatementPrint) isStatement()  {}
func (*StatementPrint) isDeferrable() {}

func (*StatementBlock) isNode()       {}
func (*StatementBlock) isStatement()  {}
func (*StatementBlock) isDeferrable() {}

func (*StatementReturn) isNode()      {}
func (*StatementReturn) isStatement() {}

func (*StatementDefer) isNode()      {}
func (*StatementDefer) isStatement() {}

func (*StatementBreak) isNode()      {}
func (*StatementBreak) isStatement() {}

func (*StatementContinue) isNode()      {}
func (*StatementContinue) isStatement() {}

func (*StatementFor) isNode()      {}
func (*StatementFor) isStatement() {}

func (*StatementWhile) isNode()      {}
func (*StatementWhile) isStatement() {}

func (*StatementIfElse) isNode()      {}
func (*StatementIfElse) isStatement() {}
