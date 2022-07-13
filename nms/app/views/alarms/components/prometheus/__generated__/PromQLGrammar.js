// @ts-nocheck
// @generated
/* eslint-disabled */
// Generated automatically by nearley, version 2.19.1
// http://github.com/Hardmath123/nearley
(function () {
function id(x) { return x[0]; }

const {lexer} = require('../PromQLTokenizer');
const {FUNCTION_NAMES, SyntaxError} = require('../PromQLTypes')
const {
        AggregationOperation,
        BinaryComparator,
        BinaryOperation,
        Clause,
        Function,
        InstantSelector,
        Label,
        Labels,
        RangeSelector,
        Scalar,
        String,
        SubQuery,
        VectorMatchClause
} = require('../PromQL');
var grammar = {
    Lexer: lexer,
    ParserRules: [
    {"name": "expression", "symbols": ["metric_selector"], "postprocess": id},
    {"name": "expression", "symbols": ["aggregation"], "postprocess": id},
    {"name": "expression", "symbols": ["function"], "postprocess": id},
    {"name": "expression", "symbols": ["binary_operation"], "postprocess": id},
    {"name": "expression", "symbols": ["subquery"], "postprocess": id},
    {"name": "expression", "symbols": ["SCALAR"], "postprocess": id},
    {"name": "metric_selector", "symbols": ["selector"], "postprocess": id},
    {"name": "metric_selector", "symbols": ["selector", "offset_clause"], "postprocess": ([selector, offset]) => selector.setOffset(offset)},
    {"name": "selector", "symbols": ["instant_selector"], "postprocess": id},
    {"name": "selector", "symbols": ["range_selector"], "postprocess": id},
    {"name": "instant_selector", "symbols": ["metric_identifier", "label_selector"], "postprocess": ([id, labels]) => new InstantSelector(id, labels)},
    {"name": "instant_selector", "symbols": ["metric_identifier"], "postprocess": ([id]) => new InstantSelector(id)},
    {"name": "instant_selector", "symbols": ["label_selector"], "postprocess": ([labels]) => new InstantSelector("", labels)},
    {"name": "metric_identifier$ebnf$1", "symbols": [(lexer.has("colon") ? {type: "colon"} : colon)]},
    {"name": "metric_identifier$ebnf$1", "symbols": ["metric_identifier$ebnf$1", (lexer.has("colon") ? {type: "colon"} : colon)], "postprocess": function arrpush(d) {return d[0].concat([d[1]]);}},
    {"name": "metric_identifier$ebnf$2", "symbols": ["IDENTIFIER"], "postprocess": id},
    {"name": "metric_identifier$ebnf$2", "symbols": [], "postprocess": function(d) {return null;}},
    {"name": "metric_identifier", "symbols": ["metric_identifier", "metric_identifier$ebnf$1", "metric_identifier$ebnf$2"], "postprocess": matches => matches.join('')},
    {"name": "metric_identifier", "symbols": [(lexer.has("colon") ? {type: "colon"} : colon), "metric_identifier"], "postprocess": matches => matches.join('')},
    {"name": "metric_identifier", "symbols": ["IDENTIFIER"], "postprocess": id},
    {"name": "range_selector", "symbols": ["instant_selector", "duration"], "postprocess": ([selector, duration]) => new RangeSelector(selector, duration)},
    {"name": "duration", "symbols": [(lexer.has("lBracket") ? {type: "lBracket"} : lBracket), "RANGE", (lexer.has("rBracket") ? {type: "rBracket"} : rBracket)], "postprocess": ([_lBracket, range, _rBracket]) => range},
    {"name": "binary_operation", "symbols": ["expression", "bin_op", "expression"], "postprocess": ([lh, op, rh]) => new BinaryOperation(lh, rh, op)},
    {"name": "binary_operation", "symbols": ["expression", "bin_op", "vector_match_clause", "expression"], "postprocess": ([lh, op, clause, rh]) => new BinaryOperation(lh, rh, op, clause)},
    {"name": "vector_match_clause", "symbols": ["MATCH_CLAUSE", "labelList"], "postprocess": ([op, labels]) => new VectorMatchClause(new Clause(op, labels))},
    {"name": "vector_match_clause", "symbols": ["MATCH_CLAUSE", "labelList", "GROUP_CLAUSE", "labelList"], "postprocess":  ([matchOp, matchLabels, groupOp, groupLabels]) =>
        new VectorMatchClause(
                new Clause(matchOp, matchLabels),
                new Clause(groupOp, groupLabels))
                                },
    {"name": "vector_match_clause", "symbols": ["MATCH_CLAUSE", "labelList", "GROUP_CLAUSE"], "postprocess":  ([matchOp, matchLabels, groupOp]) =>
        new VectorMatchClause(
                new Clause(matchOp, matchLabels),
                new Clause(groupOp))
                                },
    {"name": "bin_op", "symbols": ["BIN_COMP"], "postprocess": ([op, _boolMode]) => new BinaryComparator(op)},
    {"name": "bin_op", "symbols": ["BIN_COMP", {"literal":"bool"}], "postprocess": ([op, _boolMode]) => new BinaryComparator(op).makeBoolean()},
    {"name": "bin_op", "symbols": ["SET_OP"], "postprocess": id},
    {"name": "bin_op", "symbols": ["ARITHM_OP"], "postprocess": id},
    {"name": "offset_clause", "symbols": [{"literal":"offset"}, "RANGE"], "postprocess": ([_keyword, offset]) => offset},
    {"name": "aggregation", "symbols": ["AGG_OP", "func_call_body"], "postprocess": ([aggOp, params]) => new AggregationOperation(aggOp, params)},
    {"name": "aggregation", "symbols": ["AGG_OP", "func_call_body", "dimensionClause"], "postprocess":  ([aggOp, params, clause]) =>
        new AggregationOperation(aggOp, params, clause)
                                },
    {"name": "aggregation", "symbols": ["AGG_OP", "dimensionClause", "func_call_body"], "postprocess":  ([aggOp, clause, params,]) =>
        new AggregationOperation(aggOp, params, clause)
                                },
    {"name": "dimensionClause", "symbols": ["AGG_CLAUSE", "labelList"], "postprocess": ([op, labelList]) => new Clause(op, labelList)},
    {"name": "labelList", "symbols": [(lexer.has("lParen") ? {type: "lParen"} : lParen), "label_name_list", (lexer.has("rParen") ? {type: "rParen"} : rParen)], "postprocess": ([_lParen, labels, _rParen]) => labels},
    {"name": "label_name_list", "symbols": ["label_name_list", (lexer.has("comma") ? {type: "comma"} : comma), "IDENTIFIER"], "postprocess": ([existingLabels, _, newLabel]) => [...existingLabels, newLabel]},
    {"name": "label_name_list", "symbols": ["IDENTIFIER"], "postprocess": d => [d[0]]},
    {"name": "label_selector", "symbols": [(lexer.has("lBrace") ? {type: "lBrace"} : lBrace), "label_match_list", (lexer.has("rBrace") ? {type: "rBrace"} : rBrace)], "postprocess":  ([_lBrace, labels, _rBrace]) => {
                const ret = new Labels();
                labels.forEach(l => ret.addLabel(l.name, l.value, l.operator));
                return ret
        }},
    {"name": "label_selector", "symbols": [(lexer.has("lBrace") ? {type: "lBrace"} : lBrace), (lexer.has("rBrace") ? {type: "rBrace"} : rBrace)], "postprocess": d => new Labels()},
    {"name": "label_match_list", "symbols": ["label_match_list", (lexer.has("comma") ? {type: "comma"} : comma), "label_matcher"], "postprocess": ([existingLabels, _, newLabel]) => [...existingLabels, newLabel]},
    {"name": "label_match_list", "symbols": ["label_matcher"], "postprocess": d => [d[0]]},
    {"name": "label_matcher", "symbols": ["label", "LABEL_OP", "STRING"], "postprocess": ([name, op, value]) => new Label(name, value, op)},
    {"name": "label", "symbols": ["IDENTIFIER"], "postprocess": id},
    {"name": "label", "symbols": ["SET_OP"], "postprocess": id},
    {"name": "label", "symbols": ["AGG_CLAUSE"], "postprocess": id},
    {"name": "label", "symbols": ["MATCH_CLAUSE"], "postprocess": id},
    {"name": "label", "symbols": ["GROUP_CLAUSE"], "postprocess": id},
    {"name": "function", "symbols": ["IDENTIFIER", "func_call_body"], "postprocess":  ([funcName, params]) => {
                if (FUNCTION_NAMES.includes(funcName)) {
                        return new Function(funcName, params)
                } else {
                        throw new SyntaxError(`Unknown function: ${funcName}`);
                }
        }},
    {"name": "func_call_body", "symbols": [(lexer.has("lParen") ? {type: "lParen"} : lParen), "func_params", (lexer.has("rParen") ? {type: "rParen"} : rParen)], "postprocess": ([_lParen, params, _rParen]) => params},
    {"name": "func_params", "symbols": ["func_params", (lexer.has("comma") ? {type: "comma"} : comma), "parameter"], "postprocess": ([existingParams, _comma, newParam]) => [...existingParams, newParam]},
    {"name": "func_params", "symbols": ["parameter"], "postprocess": d => [d[0]]},
    {"name": "parameter", "symbols": ["SCALAR"], "postprocess": id},
    {"name": "parameter", "symbols": ["expression"], "postprocess": id},
    {"name": "parameter", "symbols": ["STRING"], "postprocess": d => new String(d[0])},
    {"name": "subquery$ebnf$1", "symbols": ["RANGE"], "postprocess": id},
    {"name": "subquery$ebnf$1", "symbols": [], "postprocess": function(d) {return null;}},
    {"name": "subquery$ebnf$2", "symbols": ["offset_clause"], "postprocess": id},
    {"name": "subquery$ebnf$2", "symbols": [], "postprocess": function(d) {return null;}},
    {"name": "subquery", "symbols": ["expression", (lexer.has("lBracket") ? {type: "lBracket"} : lBracket), "RANGE", (lexer.has("colon") ? {type: "colon"} : colon), "subquery$ebnf$1", (lexer.has("rBracket") ? {type: "rBracket"} : rBracket), "subquery$ebnf$2"], "postprocess":  ([expr, _lBRacket, range, _colon, step, _rBracket, offset]) =>
        new SubQuery(
                expr,
                range,
                step === null ? undefined : step,
                offset === null ? undefined : offset)
                                },
    {"name": "SCALAR", "symbols": [(lexer.has("scalar") ? {type: "scalar"} : scalar)], "postprocess": d => new Scalar(d[0].value)},
    {"name": "STRING", "symbols": [(lexer.has("string") ? {type: "string"} : string)], "postprocess": d => d[0].value},
    {"name": "IDENTIFIER", "symbols": [(lexer.has("identifier") ? {type: "identifier"} : identifier)], "postprocess": d => d[0].value},
    {"name": "LABEL_OP", "symbols": [(lexer.has("labelOp") ? {type: "labelOp"} : labelOp)], "postprocess": d => d[0].value},
    {"name": "LABEL_OP", "symbols": [(lexer.has("neq") ? {type: "neq"} : neq)], "postprocess": d => d[0].value},
    {"name": "BIN_COMP", "symbols": [(lexer.has("binComp") ? {type: "binComp"} : binComp)], "postprocess": d => d[0].value},
    {"name": "BIN_COMP", "symbols": [(lexer.has("neq") ? {type: "neq"} : neq)], "postprocess": d => d[0].value},
    {"name": "SET_OP", "symbols": [(lexer.has("setOp") ? {type: "setOp"} : setOp)], "postprocess": d => d[0].value},
    {"name": "ARITHM_OP", "symbols": [(lexer.has("arithmetic") ? {type: "arithmetic"} : arithmetic)], "postprocess": d => d[0].value},
    {"name": "AGG_OP", "symbols": [(lexer.has("aggOp") ? {type: "aggOp"} : aggOp)], "postprocess": d => d[0].value},
    {"name": "AGG_CLAUSE", "symbols": [(lexer.has("aggClause") ? {type: "aggClause"} : aggClause)], "postprocess": d => d[0].value},
    {"name": "FUNC_NAME", "symbols": [(lexer.has("functionName") ? {type: "functionName"} : functionName)], "postprocess": d => d[0].value},
    {"name": "RANGE", "symbols": [(lexer.has("range") ? {type: "range"} : range)], "postprocess": d => d[0].value},
    {"name": "MATCH_CLAUSE", "symbols": [(lexer.has("matchClause") ? {type: "matchClause"} : matchClause)], "postprocess": d => d[0].value},
    {"name": "GROUP_CLAUSE", "symbols": [(lexer.has("groupClause") ? {type: "groupClause"} : groupClause)], "postprocess": d => d[0].value}
]
  , ParserStart: "expression"
}
if (typeof module !== 'undefined'&& typeof module.exports !== 'undefined') {
   module.exports = grammar;
} else {
   window.grammar = grammar;
}
})();
