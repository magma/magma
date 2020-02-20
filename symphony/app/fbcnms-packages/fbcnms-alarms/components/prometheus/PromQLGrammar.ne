@{%
const {lexer} = require('../PromQLTokenizer');
const {FUNCTION_NAMES, SyntaxError} = require('../PromQLTypes')
const {AggregationOperation, BinaryOperation, Clause, Function, InstantSelector, Label, Labels, RangeSelector, Scalar, String, VectorMatchClause} = require('../PromQL');
%}

@lexer lexer

expression -> metric_selector  {% id %}
            | aggregation      {% id %}
            | function         {% id %}
            | binary_operation {% id %}
            | SCALAR           {% id %}

metric_selector -> selector               {% id %}
                 | selector offset_clause {% ([selector, offset]) => selector.setOffset(offset[1]) %}

selector -> instant_selector {% id %}
          | range_selector   {% id %}

instant_selector -> IDENTIFIER label_selector {% ([id, labels]) => new InstantSelector(id, labels) %}
                  | IDENTIFIER                {% ([id]) => new InstantSelector(id) %}
                  | label_selector            {% ([labels]) => new InstantSelector("", labels) %}

range_selector -> instant_selector duration {% ([selector, duration]) => new RangeSelector(selector, duration)%}

duration ->  %lBracket RANGE %rBracket {% ([_lBracket, range, _rBracket]) => range %}

binary_operation -> expression bin_op expression                     {% ([lh, op, rh]) => new BinaryOperation(lh, rh, op) %}
                  | expression bin_op vector_match_clause expression {% ([lh, op, clause, rh]) => new BinaryOperation(lh, rh, op, clause) %}

vector_match_clause -> CLAUSE_OP labelList                    {% ([op, labels]) => new VectorMatchClause(new Clause(op, labels)) %}
                     | CLAUSE_OP labelList GROUP_OP labelList {% ([matchOp, matchLabels, groupOp, groupLabels]) => new VectorMatchClause(new Clause(matchOp, matchLabels), new Clause(groupOp, groupLabels)) %}

bin_op  ->  BIN_COMP    {% id %}
        |   SET_OP      {% id %}
        |   ARITHM_OP   {% id %}

offset_clause -> "offset" RANGE

aggregation -> AGG_OP %lParen func_params %rParen                 {% ([aggOp, _lParen, params, _rParen]) => new AggregationOperation(aggOp, params) %}
             | AGG_OP %lParen func_params %rParen dimensionClause {% ([aggOp, _lParen, params, _rParen, clause]) => new AggregationOperation(aggOp, params, clause) %}
             | AGG_OP dimensionClause %lParen func_params %rParen {% ([aggOp, clause, _lParen, params, _rParen]) => new AggregationOperation(aggOp, params, clause) %}

dimensionClause -> CLAUSE_OP labelList {% ([op, labelList]) => new Clause(op, labelList) %}

labelList -> %lParen label_name_list %rParen {% ([_lParen, labels, _rParen]) => labels %}

label_name_list -> label_name_list %comma IDENTIFIER {% ([existingLabels, _, newLabel]) => [...existingLabels, newLabel] %}
                 | IDENTIFIER                        {% d => [d[0]] %}

label_selector -> %lBrace label_match_list %rBrace {% ([_lBrace, labels, _rBrace]) => {const ret = new Labels(); labels.forEach(l => ret.addLabel(l.name, l.value, l.operator)); return ret} %}
                | %lBrace %rBrace                  {% d => new Labels() %}

label_match_list -> label_match_list %comma label_matcher {% ([existingLabels, _, newLabel]) => [...existingLabels, newLabel] %}
                  | label_matcher                         {% d => [d[0]] %}

label_matcher -> label LABEL_OP STRING {% ([name, op, value]) => new Label(name, value, op) %}

label   ->  IDENTIFIER  {% id %}
        |   SET_OP      {% id %}
        |   GROUP_OP    {% id %}
        |   CLAUSE_OP   {% id %}

function -> IDENTIFIER %lParen func_params %rParen {% ([funcName, _lParen, params, _rParen]) => {
        if (FUNCTION_NAMES.includes(funcName)) {
             return new Function(funcName, params)
         } else {
             throw new SyntaxError(`Unknown function: ${funcName}`);
         }
}
%}

func_params -> func_params %comma parameter {% ([existingParams, _comma, newParam]) => [...existingParams, newParam] %}
             | parameter                    {% d => [d[0]] %}

parameter -> SCALAR     {% id %}
           | expression {% id %}
           | STRING     {% d => new String(d[0]) %}

# Terminals
SCALAR      ->  %scalar         {% d => new Scalar(d[0].value) %}
STRING      ->  %string         {% d => d[0].value %}
IDENTIFIER  ->  %identifier     {% d => d[0].value %}
LABEL_OP    ->  %labelOp        {% d => d[0].value %}
BIN_COMP    ->  %binComp        {% d => d[0].value %}
SET_OP      ->  %setOp          {% d => d[0].value %}
ARITHM_OP   ->  %arithmetic     {% d => d[0].value %}
AGG_OP      ->  %aggOp          {% d => d[0].value %}
FUNC_NAME   ->  %functionName   {% d => d[0].value %}
RANGE       ->  %range          {% d => d[0].value %}
CLAUSE_OP   ->  %clauseOp       {% d => d[0].value %}
GROUP_OP    ->  %groupOp        {% d => d[0].value %}
