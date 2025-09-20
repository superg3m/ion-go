package AST

func (s StatementAssignment) isNode()      {}
func (s StatementAssignment) isStatement() {}

func (s StatementPrint) isNode()      {}
func (s StatementPrint) isStatement() {}

func (s StatementBlock) isNode()      {}
func (s StatementBlock) isStatement() {}

func (s StatementReturn) isNode()      {}
func (s StatementReturn) isStatement() {}

func (s StatementFor) isNode()      {}
func (s StatementFor) isStatement() {}
