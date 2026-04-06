package AST

func (*ExpressionCharacter) isNode()       {}
func (*ExpressionCharacter) isExpression() {}

func (*ExpressionInteger) isNode()       {}
func (*ExpressionInteger) isExpression() {}

func (*ExpressionFloat) isNode()       {}
func (*ExpressionFloat) isExpression() {}

func (*ExpressionBoolean) isNode()       {}
func (*ExpressionBoolean) isExpression() {}

func (*ExpressionString) isNode()       {}
func (*ExpressionString) isExpression() {}

func (*ExpressionIdentifier) isNode()       {}
func (*ExpressionIdentifier) isExpression() {}

func (*ExpressionGrouping) isNode()       {}
func (*ExpressionGrouping) isExpression() {}

func (*ExpressionTypeCast) isNode()       {}
func (*ExpressionTypeCast) isExpression() {}

func (*ExpressionUnary) isNode()       {}
func (*ExpressionUnary) isExpression() {}

func (*ExpressionBinary) isNode()       {}
func (*ExpressionBinary) isExpression() {}

func (*ExpressionArray) isNode()       {}
func (*ExpressionArray) isExpression() {}

func (*ExpressionArrayAccess) isNode()       {}
func (*ExpressionArrayAccess) isExpression() {}

func (*ExpressionStruct) isNode()       {}
func (*ExpressionStruct) isExpression() {}

func (*ExpressionAccessChain) isNode()       {}
func (*ExpressionAccessChain) isExpression() {}

func (*ExpressionLen) isNode()       {}
func (*ExpressionLen) isExpression() {}

func (*ExpressionPseudo) isNode()       {}
func (*ExpressionPseudo) isExpression() {}
