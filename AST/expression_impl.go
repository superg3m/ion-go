package AST

func (*ExpressionInteger) isNode()       {}
func (*ExpressionInteger) isExpression() {}

func (*ExpressionFloat) isNode()       {}
func (*ExpressionFloat) isExpression() {}

func (*ExpressionBoolean) isNode()       {}
func (*ExpressionBoolean) isExpression() {}

func (e ExpressionString) isNode()       {}
func (e ExpressionString) isExpression() {}

func (*ExpressionIdentifier) isNode()       {}
func (*ExpressionIdentifier) isExpression() {}

func (*ExpressionGrouping) isNode()       {}
func (*ExpressionGrouping) isExpression() {}

func (*ExpressionUnary) isNode()       {}
func (*ExpressionUnary) isExpression() {}

func (*ExpressionBinary) isNode()       {}
func (*ExpressionBinary) isExpression() {}

func (*ExpressionFunctionCall) isNode()       {}
func (*ExpressionFunctionCall) isExpression() {}
func (*ExpressionFunctionCall) isDeferrable() {}

func (*ExpressionArray) isNode()       {}
func (*ExpressionArray) isExpression() {}

func (*ExpressionLen) isNode()       {}
func (*ExpressionLen) isExpression() {}

func (*ExpressionArrayAccess) isNode()       {}
func (*ExpressionArrayAccess) isExpression() {}

func (*ExpressionPseudo) isNode()       {}
func (*ExpressionPseudo) isExpression() {}
